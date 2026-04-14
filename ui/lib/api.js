const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

async function request(path, options = {}) {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {})
    }
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`API request failed (${res.status}): ${text}`);
  }

  if (res.status === 204) {
    return null;
  }

  return res.json();
}

export async function getHealth() {
  return request("/health");
}

export async function getProjects() {
  return request("/projects");
}

export async function createProject(payload) {
  return request("/projects", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function getProject(id) {
  return request(`/projects/${id}`);
}

export async function getProjectTasks(projectId) {
  return request(`/projects/${projectId}/tasks`);
}

export async function createTask(projectId, payload) {
  return request(`/projects/${projectId}/tasks`, {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function updateProject(id, payload) {
  return request(`/projects/${id}`, {
    method: "PUT",
    body: JSON.stringify(payload)
  });
}

export async function deleteProject(id) {
  return request(`/projects/${id}`, {
    method: "DELETE"
  });
}

export async function updateTask(id, payload) {
  return request(`/tasks/${id}`, {
    method: "PUT",
    body: JSON.stringify(payload)
  });
}

export async function deleteTask(id) {
  return request(`/tasks/${id}`, {
    method: "DELETE"
  });
}
