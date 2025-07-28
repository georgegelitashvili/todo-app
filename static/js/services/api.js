
export class API {
    static getAuthHeader() {
        const token = localStorage.getItem('token');
        console.log('Current token:', token); // Debug log
        console.log('Token type:', typeof token); // Debug log
        console.log('Token null check:', token === null, token === 'null'); // Debug log
        
        if (!token || token === 'null') {
            throw new Error('No authentication token found');
        }
        return `Bearer ${token}`;
    }

    static async login(credentials) {
        const response = await fetch('/api/v1/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(credentials)
        });
        return await response.json();
    }

    static async register(userData) {
        const response = await fetch('/api/v1/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(userData)
        });
        return await response.json();
    }

    static async getTasks() {
        try {
            const response = await fetch('/api/v1/tasks', {
                headers: {
                    'Authorization': this.getAuthHeader()
                }
            });
            
            if (response.status === 401) {
                console.log('Token expired or invalid, redirecting to login');
                localStorage.removeItem('token');
                window.location.href = '/login';
                return { status: 'error', message: 'Unauthorized' };
            }
            
            return await response.json();
        } catch (error) {
            console.error('Error in getTasks:', error);
            
            // Check if it's an auth error
            if (error.message && error.message.includes('authentication token')) {
                console.log('No valid token, redirecting to login');
                localStorage.removeItem('token');
                window.location.href = '/login';
                return { status: 'error', message: 'Unauthorized' };
            }
            
            throw error;
        }
    }

    static async createTask(taskData) {
        try {
            const response = await fetch('/api/v1/tasks', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': this.getAuthHeader()
                },
                body: JSON.stringify(taskData)
            });
            
            if (response.status === 401) {
                console.log('Token expired or invalid, redirecting to login');
                localStorage.removeItem('token');
                window.location.href = '/login';
                return { status: 'error', message: 'Unauthorized' };
            }
            
            return await response.json();
        } catch (error) {
            console.error('Error in createTask:', error);
            
            // Check if it's an auth error
            if (error.message && error.message.includes('authentication token')) {
                console.log('No valid token, redirecting to login');
                localStorage.removeItem('token');
                window.location.href = '/login';
                return { status: 'error', message: 'Unauthorized' };
            }
            
            throw error;
        }
    }


    static async deleteTask(taskId) {
        try {
            const response = await fetch(`/api/v1/tasks/${taskId}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': this.getAuthHeader()
                }
            });
            
            if (response.status === 401) {
                console.log('Token expired or invalid, redirecting to login');
                localStorage.removeItem('token');
                window.location.href = '/login';
                return { status: 'error', message: 'Unauthorized' };
            }
            
            return await response.json();
        } catch (error) {
            console.error('Error in deleteTask:', error);
            
            // Check if it's an auth error
            if (error.message && error.message.includes('authentication token')) {
                console.log('No valid token, redirecting to login');
                localStorage.removeItem('token');
                window.location.href = '/login';
                return { status: 'error', message: 'Unauthorized' };
            }
            
            throw error;
        }
    }

    static async updateTask(taskId, taskData) {
        try {
            const response = await fetch(`/api/v1/tasks/${taskId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': this.getAuthHeader()
                },
                body: JSON.stringify(taskData)
            });
            
            if (response.status === 401) {
                console.log('Token expired or invalid, redirecting to login');
                localStorage.removeItem('token');
                window.location.href = '/login';
                return { status: 'error', message: 'Unauthorized' };
            }
            
            return await response.json();
        } catch (error) {
            console.error('Error in updateTask:', error);
            
            // Check if it's an auth error
            if (error.message && error.message.includes('authentication token')) {
                console.log('No valid token, redirecting to login');
                localStorage.removeItem('token');
                window.location.href = '/login';
                return { status: 'error', message: 'Unauthorized' };
            }
            
            throw error;
        }
    }
}
