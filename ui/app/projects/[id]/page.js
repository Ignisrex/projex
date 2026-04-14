"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { createTask, getProject, getProjectTasks } from "../../../lib/api";

export default function ProjectDetailsPage({ params }) {
  const [project, setProject] = useState(null);
  const [tasks, setTasks] = useState([]);
  const [error, setError] = useState("");
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [status, setStatus] = useState("todo");
  const [priority, setPriority] = useState(1);

  async function load() {
    setError("");
    try {
      const [projectData, tasksData] = await Promise.all([
        getProject(params.id),
        getProjectTasks(params.id)
      ]);
      setProject(projectData);
      setTasks(tasksData || []);
    } catch (e) {
      setProject(null);
      setTasks([]);
      setError(e.message);
    }
  }

  useEffect(() => {
    load();
  }, [params.id]);

  async function onCreateTask(e) {
    e.preventDefault();
    setError("");
    try {
      await createTask(params.id, {
        title,
        description,
        status,
        priority: Number(priority)
      });
      setTitle("");
      setDescription("");
      setStatus("todo");
      setPriority(1);
      await load();
    } catch (e) {
      setError(e.message);
    }
  }

  return (
    <main className="container">
      <p><Link href="/">← Back to projects</Link></p>
      <h1>Project</h1>
      {error ? <p style={{ color: "#b91c1c" }}>{error}</p> : null}

      <section className="card">
        {project ? (
          <>
            <h2>{project.name || "Unnamed Project"}</h2>
            <p>{project.description || "No description"}</p>
          </>
        ) : (
          <p>Project not available yet (API scaffold).</p>
        )}
      </section>

      <section className="card">
        <h2>Tasks</h2>
        {tasks.length === 0 ? (
          <p>No tasks yet (or API methods are still scaffolded).</p>
        ) : (
          <ul>
            {tasks.map((task) => (
              <li key={task.id}>
                {task.title} — {task.status} (priority: {task.priority})
              </li>
            ))}
          </ul>
        )}
      </section>

      <section className="card">
        <h2>Create Task</h2>
        <form onSubmit={onCreateTask}>
          <div className="row">
            <div style={{ flex: 2, minWidth: "220px" }}>
              <label htmlFor="title">Title</label>
              <input id="title" value={title} onChange={(e) => setTitle(e.target.value)} required />
            </div>
            <div style={{ flex: 2, minWidth: "220px" }}>
              <label htmlFor="taskDescription">Description</label>
              <input
                id="taskDescription"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
              />
            </div>
          </div>
          <div className="row" style={{ marginTop: "0.75rem" }}>
            <div style={{ flex: 1, minWidth: "180px" }}>
              <label htmlFor="status">Status</label>
              <select id="status" value={status} onChange={(e) => setStatus(e.target.value)}>
                <option value="todo">todo</option>
                <option value="in_progress">in_progress</option>
                <option value="done">done</option>
              </select>
            </div>
            <div style={{ flex: 1, minWidth: "180px" }}>
              <label htmlFor="priority">Priority</label>
              <input
                id="priority"
                type="number"
                min="1"
                max="5"
                value={priority}
                onChange={(e) => setPriority(e.target.value)}
              />
            </div>
          </div>
          <div style={{ marginTop: "0.75rem" }}>
            <button type="submit">Create Task</button>
          </div>
        </form>
      </section>
    </main>
  );
}
