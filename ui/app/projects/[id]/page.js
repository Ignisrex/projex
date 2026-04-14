"use client";

import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import {
  getProject,
  getProjectTasks,
  createTask,
  updateProject,
  deleteProject,
  updateTask,
  deleteTask
} from "../../../lib/api";

export default function ProjectDetailsPage() {
  const { id } = useParams();
  const router = useRouter();

  const [project, setProject] = useState(null);
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  // Project edit state
  const [editingProject, setEditingProject] = useState(false);
  const [editProjectName, setEditProjectName] = useState("");
  const [editProjectDescription, setEditProjectDescription] = useState("");

  // New task form state
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [status, setStatus] = useState("todo");
  const [priority, setPriority] = useState(1);

  // Task edit state
  const [editingTaskId, setEditingTaskId] = useState(null);
  const [editTaskTitle, setEditTaskTitle] = useState("");
  const [editTaskDescription, setEditTaskDescription] = useState("");
  const [editTaskStatus, setEditTaskStatus] = useState("todo");
  const [editTaskPriority, setEditTaskPriority] = useState(1);

  async function load() {
    setLoading(true);
    setError("");
    try {
      const [projectData, tasksData] = await Promise.all([
        getProject(id),
        getProjectTasks(id)
      ]);
      setProject(projectData);
      setEditProjectName(projectData.name || "");
      setEditProjectDescription(projectData.description || "");
      setTasks(tasksData || []);
    } catch (e) {
      setProject(null);
      setTasks([]);
      setError(e.message);
    }
    setLoading(false);
  }

  async function loadTasks() {
    setError("");
    try {
      const data = await getProjectTasks(id);
      setTasks(data || []);
    } catch (e) {
      setError(e.message);
    }
  }

  useEffect(() => {
    if (id) load();
  }, [id]);

  // --- Create Task ---
  async function onCreateTask(e) {
    e.preventDefault();
    setError("");
    try {
      await createTask(id, {
        title,
        description,
        status,
        priority: Number(priority)
      });
      setTitle("");
      setDescription("");
      setStatus("todo");
      setPriority(1);
      await loadTasks();
    } catch (e) {
      setError(e.message);
    }
  }

  // --- Delete Task ---
  async function onDeleteTask(taskId) {
    setError("");
    try {
      await deleteTask(taskId);
      await loadTasks();
    } catch (e) {
      setError(e.message);
    }
  }

  // --- Edit Task ---
  function startEditingTask(task) {
    setEditingTaskId(task.id);
    setEditTaskTitle(task.title || "");
    setEditTaskDescription(task.description || "");
    setEditTaskStatus(task.status || "todo");
    setEditTaskPriority(task.priority ?? 1);
  }

  function cancelEditingTask() {
    setEditingTaskId(null);
  }

  async function onUpdateTask(e) {
    e.preventDefault();
    setError("");
    try {
      await updateTask(editingTaskId, {
        title: editTaskTitle,
        description: editTaskDescription,
        status: editTaskStatus,
        priority: Number(editTaskPriority)
      });
      setEditingTaskId(null);
      await loadTasks();
    } catch (e) {
      setError(e.message);
    }
  }

  // --- Edit Project ---
  async function onUpdateProject(e) {
    e.preventDefault();
    setError("");
    try {
      const updated = await updateProject(id, {
        name: editProjectName,
        description: editProjectDescription
      });
      setProject(updated);
      setEditingProject(false);
    } catch (e) {
      setError(e.message);
    }
  }

  // --- Delete Project ---
  async function onDeleteProject() {
    setError("");
    try {
      await deleteProject(id);
      router.push("/");
    } catch (e) {
      setError(e.message);
    }
  }

  if (loading) {
    return (
      <main className="container">
        <p>Loading…</p>
      </main>
    );
  }

  return (
    <main className="container">
      <p><Link href="/">← Back to projects</Link></p>

      {error ? <p style={{ color: "#b91c1c" }}>{error}</p> : null}

      {/* --- Project Header --- */}
      <section className="card">
        {editingProject ? (
          <form onSubmit={onUpdateProject}>
            <div className="row">
              <div style={{ flex: 1, minWidth: "260px" }}>
                <label htmlFor="editProjectName">Name</label>
                <input
                  id="editProjectName"
                  value={editProjectName}
                  onChange={(e) => setEditProjectName(e.target.value)}
                  required
                />
              </div>
              <div style={{ flex: 2, minWidth: "260px" }}>
                <label htmlFor="editProjectDescription">Description</label>
                <input
                  id="editProjectDescription"
                  value={editProjectDescription}
                  onChange={(e) => setEditProjectDescription(e.target.value)}
                />
              </div>
            </div>
            <div style={{ marginTop: "0.75rem", display: "flex", gap: "0.5rem" }}>
              <button type="submit">Save</button>
              <button type="button" onClick={() => setEditingProject(false)}>
                Cancel
              </button>
            </div>
          </form>
        ) : (
          <>
            <h2>{project?.name || "Unnamed Project"}</h2>
            <p>{project?.description || "No description"}</p>
            <div style={{ display: "flex", gap: "0.5rem", marginTop: "0.75rem" }}>
              <button onClick={() => setEditingProject(true)}>Edit Project</button>
              <button
                onClick={onDeleteProject}
                style={{ background: "#b91c1c" }}
              >
                Delete Project
              </button>
            </div>
          </>
        )}
      </section>

      {/* --- Task List --- */}
      <section className="card">
        <h2>Tasks</h2>
        {tasks.length === 0 ? (
          <p>No tasks yet.</p>
        ) : (
          <ul style={{ listStyle: "none", padding: 0 }}>
            {tasks.map((task) => (
              <li
                key={task.id}
                style={{
                  borderBottom: "1px solid #e5e7eb",
                  padding: "0.75rem 0"
                }}
              >
                {editingTaskId === task.id ? (
                  <form onSubmit={onUpdateTask}>
                    <div className="row">
                      <div style={{ flex: 2, minWidth: "200px" }}>
                        <label>Title</label>
                        <input
                          value={editTaskTitle}
                          onChange={(e) => setEditTaskTitle(e.target.value)}
                          required
                        />
                      </div>
                      <div style={{ flex: 2, minWidth: "200px" }}>
                        <label>Description</label>
                        <input
                          value={editTaskDescription}
                          onChange={(e) => setEditTaskDescription(e.target.value)}
                        />
                      </div>
                      <div style={{ flex: 1, minWidth: "120px" }}>
                        <label>Status</label>
                        <select
                          value={editTaskStatus}
                          onChange={(e) => setEditTaskStatus(e.target.value)}
                        >
                          <option value="todo">todo</option>
                          <option value="in_progress">in_progress</option>
                          <option value="done">done</option>
                        </select>
                      </div>
                      <div style={{ flex: 0, minWidth: "80px" }}>
                        <label>Priority</label>
                        <input
                          type="number"
                          min="1"
                          value={editTaskPriority}
                          onChange={(e) => setEditTaskPriority(e.target.value)}
                        />
                      </div>
                    </div>
                    <div style={{ marginTop: "0.5rem", display: "flex", gap: "0.5rem" }}>
                      <button type="submit">Save</button>
                      <button type="button" onClick={cancelEditingTask}>
                        Cancel
                      </button>
                    </div>
                  </form>
                ) : (
                  <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                    <div>
                      <strong>{task.title}</strong>
                      {task.description ? <span> — {task.description}</span> : null}
                      <br />
                      <small>
                        Status: {task.status} | Priority: {task.priority}
                      </small>
                    </div>
                    <div style={{ display: "flex", gap: "0.5rem", flexShrink: 0 }}>
                      <button onClick={() => startEditingTask(task)}>Edit</button>
                      <button
                        onClick={() => onDeleteTask(task.id)}
                        style={{ background: "#b91c1c" }}
                      >
                        Delete
                      </button>
                    </div>
                  </div>
                )}
              </li>
            ))}
          </ul>
        )}
      </section>

      {/* --- Create Task Form --- */}
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
                value={priority}
                onChange={(e) => setPriority(e.target.value)}
              />
            </div>
          </div>
          <div style={{ marginTop: "0.75rem" }}>
            <button type="submit">Add Task</button>
          </div>
        </form>
      </section>
    </main>
  );
}
