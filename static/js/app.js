import {API} from './services/api.js';

class App {
    static isAuthenticated() {
        const token = localStorage.getItem('token');
        // Check for null, 'null', undefined, 'undefined', or empty string
        return !!(token && token !== 'null' && token !== 'undefined' && token.trim() !== '');
    }

    static checkAuth() {
        const currentPath = window.location.pathname;
        const isLoginPage = currentPath === '/login';
        const isRegisterPage = currentPath === '/register';
        
        // Clean up invalid tokens
        const token = localStorage.getItem('token');
        if (token === 'null' || token === 'undefined' || token === '') {
            localStorage.removeItem('token');
        }

        if (App.isAuthenticated()) {
            if (isLoginPage || isRegisterPage || currentPath === '/') {
                window.location.href = '/tasks';
            }
        } else {
            if (!isLoginPage && !isRegisterPage) {
                window.location.href = '/login';
            }
        }
    }

    static async handleLogin(event) {
        event.preventDefault();
        const form = event.target;
        const formData = new FormData(form);
        const errorContainer = document.createElement('div');
        errorContainer.className = 'error';

        try {
            const response = await fetch('/api/v1/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    email: formData.get('email'),
                    password: formData.get('password')
                })
            });

            const responseText = await response.text();
            console.log('Login response text:', responseText);

            const data = JSON.parse(responseText);
            console.log('Login response data:', data); // Debug log
            console.log('Token from response:', data.data?.token); // Debug log

            if (data && data.status === 'success' && data.data && data.data.token) {
                localStorage.setItem('token', data.data.token);
                console.log('Token stored in localStorage:', localStorage.getItem('token')); // Verify storage
                console.log('Token length:', data.data.token.length); // Check token length
                window.location.href = '/tasks';
            } else {
                console.error('Invalid response structure:', data); // Debug log
                throw new Error('Login failed: Invalid response format');
            }
        } catch (error) {
            console.error('Login error:', error);
            if (error) {
                errorContainer.textContent = 'Login failed. Please try again.';
                form.prepend(errorContainer);
            }
        }
    }

    static async handleRegister(event) {
        event.preventDefault();
        const form = event.target;
        const formData = new FormData(form);
        const errorContainer = document.createElement('div');
        errorContainer.className = 'error';

        try {
            const response = await API.register({
                username: formData.get('username'),
                email: formData.get('email'),
                password: formData.get('password')
            });

            console.log('Registration response:', response); // Debug log

            if (response.status === 'success') {
                console.log('Registration successful, redirecting to login');
                window.location.href = '/login';
            } else {
                throw new Error('Registration failed: ' + (response.message || 'Unknown error'));
            }
        } catch (error) {
            console.error('Registration error:', error);
            errorContainer.textContent = 'Registration failed. Please try again.';
            form.prepend(errorContainer);
        }
    }

    static async handleLogout(event) {
        if (event) {
            event.preventDefault();
        }

        try {
            localStorage.removeItem('token');
            window.location.href = '/login';
        } catch (error) {
            console.error('Logout error:', error);
        }
    }

    static async loadTasks() {
        try {
            const tasksList = document.getElementById('tasks-list');
            if (!tasksList) {
                console.log('Tasks list element not found - probably not on tasks page');
                return;
            }

            const response = await API.getTasks();
            console.log('API Response:', response); // Debug log

            if (response.status === 'success') {
                const tasks = Array.isArray(response.data) ? response.data : [];
                console.log('Tasks array:', tasks); // Debug log

                tasksList.innerHTML = tasks.map(task => `
                <div class="task-item" data-id="${task.task_id}">
                    <div class="task-content">
                        <h3 class="task-title">${task.title}</h3>
                        <p class="task-description">${task.description || 'No description'}</p>
                        <span class="task-status status-${task.status.toLowerCase()}">${task.status}</span>
                    </div>
                    <div class="task-actions">
                        <button onclick="App.handleEditTask('${task.task_id}')" class="btn-edit">Edit</button>
                        <button onclick="App.handleDeleteTask('${task.task_id}')" class="btn-delete">Delete</button>
                    </div>
                </div>
            `).join('');
            } else {
                console.error('Invalid response:', response);
            }
        } catch (error) {
            console.error('Error loading tasks:', error);
        }
    }

    static async handleAddTask(event) {
        event.preventDefault();
        event.stopPropagation();
        
        const form = event.target;
        const submitButton = form.querySelector('button[type="submit"]');
        const editId = submitButton.getAttribute('data-edit-id');
        
        console.log('=== HANDLE ADD TASK DEBUG ===');
        console.log('Token before task creation:', localStorage.getItem('token'));
        console.log('Is authenticated before task creation:', this.isAuthenticated());
        console.log('============================');
        
        const taskData = {
            title: form.querySelector('#new-task').value.trim(),
            description: form.querySelector('#task-description').value.trim(),
            status: form.querySelector('#task-status').value
        };

        console.log('Task data:', taskData); // Debug log
        if (!taskData.title) return false;

        try {
            if (editId) {
                // Update existing task
                await App.handleUpdateTask(editId, taskData);
            } else {
                // Create new task
                const response = await API.createTask(taskData);
                if (response.status === 'success') {
                    form.reset();
                    await App.loadTasks();
                }
            }
        } catch (error) {
            console.error('Error handling task:', error);
        }
        
        return false; // Prevent form submission
    }

    static async handleDeleteTask(taskId) {
        try {
            const response = await API.deleteTask(taskId);
            if (response.status === 'success') {
                await App.loadTasks();
            }
        } catch (error) {
            console.error('Error deleting task:', error);
        }
    }

    static async handleEditTask(taskId) {
        try {
            // Get current task data
            const tasks = document.querySelectorAll('.task-item');
            const taskElement = document.querySelector(`[data-id="${taskId}"]`);
            
            if (!taskElement) return;
            
            const title = taskElement.querySelector('.task-title').textContent;
            const description = taskElement.querySelector('.task-description').textContent;
            const status = taskElement.querySelector('.task-status').textContent;
            
            // Fill form with current data
            document.getElementById('new-task').value = title;
            document.getElementById('task-description').value = description === 'No description' ? '' : description;
            document.getElementById('task-status').value = status;
            
            // Change form to update mode
            const form = document.getElementById('add-task-form');
            const submitButton = form.querySelector('button[type="submit"]');
            submitButton.textContent = 'Update Task';
            submitButton.setAttribute('data-edit-id', taskId);
            
            // Scroll to form
            form.scrollIntoView({ behavior: 'smooth' });
            
        } catch (error) {
            console.error('Error editing task:', error);
        }
    }

    static async handleUpdateTask(taskId, taskData) {
        try {
            const response = await API.updateTask(taskId, taskData);
            if (response.status === 'success') {
                await App.loadTasks();
                // Reset form to add mode
                const form = document.getElementById('add-task-form');
                const submitButton = form.querySelector('button[type="submit"]');
                submitButton.textContent = 'Add Task';
                submitButton.removeAttribute('data-edit-id');
                form.reset();
            }
        } catch (error) {
            console.error('Error updating task:', error);
        }
    }


    static async init() {
        // Debug: Check localStorage contents
        console.log('=== DEBUG INFO ===');
        console.log('Current token in localStorage:', localStorage.getItem('token'));
        console.log('Is authenticated:', this.isAuthenticated());
        console.log('Current URL:', window.location.pathname);
        console.log('=================');
        
        // Clean up invalid tokens first
        const token = localStorage.getItem('token');
        if (token === 'null' || token === 'undefined' || token === '') {
            console.log('Cleaning up invalid token:', token);
            localStorage.removeItem('token');
        }
        
        // Only load tasks if the tasks list element exists (meaning we're on the tasks page)
        const tasksList = document.getElementById('tasks-list');
        if (this.isAuthenticated() && tasksList) {
            await this.loadTasks();
        }
    }
}

// Initialize auth check on page load
document.addEventListener('DOMContentLoaded', () => {
    App.init();
    App.checkAuth();
});

window.App = App;