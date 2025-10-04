import { Router } from "express";
import { redis } from "../redis.ts";

const router = Router();

// Create a group
router.post("/", async (req, res) => {
  const { members } = req.body;
  const groupId = await redis.incr("group:id:counter");

  if (!members || !Array.isArray(members) || members.length === 0) {
    return res.status(400).send("Members array required");
  }

  await redis.sadd(`group:${groupId}:members`, ...members);
  for (const uid of members) {
    await redis.sadd(`user:${uid}:groups`, groupId.toString());
  }

  res.json({ groupId });
});

// Add members to group
router.post("/:id/members", async (req, res) => {
  const groupId = req.params.id;
  const { members } = req.body;

  await redis.sadd(`group:${groupId}:members`, ...members);
  for (const uid of members) {
    await redis.sadd(`user:${uid}:groups`, groupId.toString());
  }

  res.send("Members added");
});

// Remove members from group
router.delete("/:id/members", async (req, res) => {
  const groupId = req.params.id;
  const { members } = req.body;

  await redis.srem(`group:${groupId}:members`, ...members);
  for (const uid of members) {
    await redis.srem(`user:${uid}:groups`, groupId.toString());
  }

  res.send("Members removed");
});

export default router;
