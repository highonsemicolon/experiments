import express from "express";
import bodyParser from "body-parser";
import notesRouter from "./routes/notes.ts";
import groupsRouter from "./routes/groups.ts";

const app = express();
app.use(bodyParser.json());

app.use((req, res, next) => {
  const userId = req.header("user-id");
  if (!userId) return res.status(401).send("Missing user-id header");
  (req as any).userId = userId;
  next();
});

app.use("/notes", notesRouter);
app.use("/groups", groupsRouter);

app.listen(3000, () => console.log("Server running on http://localhost:3000"));
