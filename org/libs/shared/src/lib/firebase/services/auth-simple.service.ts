import { Injectable } from '@angular/core';
import { Auth, signInWithEmailAndPassword, signOut, User, authState, createUserWithEmailAndPassword } from '@angular/fire/auth';
import { Observable, BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class SimpleAuthService {
  private currentUserSubject = new BehaviorSubject<User | null>(null);
  public currentUser$: Observable<User | null> = this.currentUserSubject.asObservable();

  constructor(private auth: Auth) {
    // Suscribirse a cambios en el estado de autenticación
    authState(this.auth).subscribe((user) => {
      this.currentUserSubject.next(user);
      console.log('Usuario autenticado:', user);
    });
  }

  /**
   * Inicia sesión con email y contraseña
   */
  async login(email: string, password: string): Promise<User> {
    try {
      const credential = await signInWithEmailAndPassword(this.auth, email, password);
      return credential.user;
    } catch (error: any) {
      console.error('Error en login:', error);
      throw this.handleFirebaseError(error);
    }
  }

  /**
   * Registra un nuevo usuario
   */
  async register(email: string, password: string): Promise<User> {
    try {
      const credential = await createUserWithEmailAndPassword(this.auth, email, password);
      return credential.user;
    } catch (error: any) {
      console.error('Error en registro:', error);
      throw this.handleFirebaseError(error);
    }
  }

  /**
   * Cierra sesión
   */
  async logout(): Promise<void> {
    try {
      await signOut(this.auth);
    } catch (error: any) {
      console.error('Error en logout:', error);
      throw error;
    }
  }

  /**
   * Obtiene el usuario actual
   */
  getCurrentUser(): User | null {
    return this.currentUserSubject.value;
  }

  /**
   * Verifica si el usuario está autenticado
   */
  isAuthenticated(): boolean {
    return this.getCurrentUser() !== null;
  }

  /**
   * Obtiene el token de autenticación
   */
  async getIdToken(): Promise<string | null> {
    const user = this.getCurrentUser();
    if (user) {
      return await user.getIdToken();
    }
    return null;
  }

  /**
   * Maneja errores de Firebase y los convierte en mensajes legibles
   */
  private handleFirebaseError(error: any): Error {
    let message = 'Error desconocido';
    
    switch (error.code) {
      case 'auth/user-not-found':
        message = 'Usuario no encontrado';
        break;
      case 'auth/wrong-password':
        message = 'Contraseña incorrecta';
        break;
      case 'auth/invalid-email':
        message = 'Email inválido';
        break;
      case 'auth/email-already-in-use':
        message = 'El email ya está en uso';
        break;
      case 'auth/weak-password':
        message = 'La contraseña es muy débil';
        break;
      case 'auth/invalid-credential':
        message = 'Credenciales inválidas';
        break;
      case 'auth/too-many-requests':
        message = 'Demasiados intentos fallidos. Intenta más tarde';
        break;
      default:
        message = error.message || 'Error de autenticación';
    }
    
    return new Error(message);
  }
}
