"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { createProject, deleteProject, getHealth, getProjects } from "../lib/api";

export default function ProjectsPage() {
  const [health, setHealth] = useState("checking");
  const [projects, setProjects] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");

  async function load() {
    setLoading(true);
    setError("");
    try {
      const healthData = await getHealth();
      setHealth(healthData.status || "unknown");
    } catch {
      setHealth("unreachable");
    }

    try {
      const projectsData = await getProjects();
      setProjects(projectsData || []);
    } catch (e) {
      setProjects([]);
      setError(e.message);
    }
    setLoading(false);
  }

  useEffect(() => {
    load();
  }, []);

  async function onDeleteProject(projectId) {
    setError("");
    try {
      await deleteProject(projectId);
      await load();
    } catch (e) {
      setError(e.message);
    }
  }

  async function onCreateProject(e) {
    e.preventDefault();
    setError("");
    try {
      await createProject({ name, description });
      setName("");
      setDescription("");
      await load();
    } catch (e) {
      setError(e.message);
    }
  }

  return (
    <main className="container">
      <h1>Project Tracker</h1>
      <p>Backend health: <strong>{health}</strong></p>
      {error ? <p style={{ color: "#b91c1c" }}>{error}</p> : null}

      <section className="card">
        <h2>Projects</h2>
        {loading ? (
          <p>Loading…</p>
        ) : projects.length === 0 ? (
          <p>No projects yet (or API methods are still scaffolded).</p>
        ) : (
          <ul>
            {projects.map((project) => (
              <li key={project.id} style={{ display: "flex", alignItems: "center", gap: "0.5rem" }}>
                <Link href={`/projects/${project.id}`}>{project.name || project.id}</Link>
                <button
                  onClick={() => onDeleteProject(project.id)}
                  style={{ background: "#b91c1c", fontSize: "0.8rem", padding: "0.25rem 0.5rem" }}
                >
                  Delete
                </button>
              </li>
            ))}
          </ul>
        )}
      </section>

      <section className="card">
        <h2>Create Project</h2>
        <form onSubmit={onCreateProject}>
          <div className="row">
            <div style={{ flex: 1, minWidth: "260px" }}>
              <label htmlFor="name">Name</label>
              <input id="name" value={name} onChange={(e) => setName(e.target.value)} required />
            </div>
            <div style={{ flex: 2, minWidth: "260px" }}>
              <label htmlFor="description">Description</label>
              <input id="description" value={description} onChange={(e) => setDescription(e.target.value)} />
            </div>
          </div>
          <div style={{ marginTop: "0.75rem" }}>
            <button type="submit">Create</button>
          </div>
        </form>
      </section>
    </main>
  );
}
