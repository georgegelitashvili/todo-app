class AuthService {
    static getToken() {
        const token = localStorage.getItem('token');
        if (!token || token === 'null' || token === 'undefined' || token.trim() === '') {
            return null;
        }
        return token;
    }

    static setToken(token) {
        if (token && token !== 'null' && token !== 'undefined') {
            localStorage.setItem('token', token);
        }
    }

    static removeToken() {
        localStorage.removeItem('token');
    }

    static isAuthenticated() {
        return !!this.getToken();
    }
}