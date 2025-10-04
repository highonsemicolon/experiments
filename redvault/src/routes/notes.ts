import { Router } from "express";
import { redis } from "../redis.ts";

const router = Router();

// Create Note
router.post("/", async (req, res) => {
  const userId = (req as any).userId;
  const { title, content, type, groupId } = req.body;

  const noteId = await redis.incr("note:id:counter");
  const now = new Date().toISOString();

  await redis.hset(`note:${noteId}`, {
    title,
    content,
    owner: userId,
    type,
    group_id: groupId || "",
    created_at: now,
    updated_at: now,
  });

  await redis.zadd(`user:${userId}:notes`, Date.now(), noteId.toString());

  // If public, add to public index
  if (type === "public") {
    await redis.zadd("notes:public", Date.now(), noteId.toString());
  }

  // If group, add to group notes set
  if (type === "group" && groupId) {
    await redis.sadd(`group:${groupId}:notes`, noteId.toString());
  }

  res.json({ noteId });
});

// Get all notes user can access
router.get("/", async (req, res) => {
  const userId = (req as any).userId;

  // 1. Owned
  const owned = await redis.zrange(`user:${userId}:notes`, 0, -1);

  // 2. Shared
  const shared = await redis.smembers(`user:${userId}:shared_notes`);

  // 3. Public
  const publicNotes = await redis.zrange("notes:public", 0, -1);

  // 4. Group notes
  const groupIds = await redis.smembers(`user:${userId}:groups`);
  let groupNotes: string[] = [];
  for (const gid of groupIds) {
    const notes = await redis.smembers(`group:${gid}:notes`);
    groupNotes = groupNotes.concat(notes);
  }

  const allNotes = Array.from(new Set([...owned, ...shared, ...publicNotes, ...groupNotes]));

  // Fetch note hashes
  const pipeline = redis.pipeline();
  allNotes.forEach((id) => pipeline.hgetall(`note:${id}`));
  const execResult = await pipeline.exec();
  const notesData = execResult ? execResult.map(([err, data]) => data) : [];

  res.json(notesData);
});

// Get a single note
router.get("/:id", async (req, res) => {
  const userId = (req as any).userId;
  const noteId = req.params.id;

  const note = await redis.hgetall(`note:${noteId}`);
  if (!note || Object.keys(note).length === 0) return res.status(404).send("Note not found");

  // Access control
  if (note.type === "public" || note.owner === userId) {
    return res.json(note);
  } else if (note.type === "group") {
    const isMember = await redis.sismember(`group:${note.group_id}:members`, userId);
    if (isMember) return res.json(note);
  }
  return res.status(403).send("Forbidden");
});

export default router;
