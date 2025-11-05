import { Injectable, Injector, runInInjectionContext } from '@angular/core';
import { Auth, signInWithEmailAndPassword, signOut, User, authState, getIdTokenResult } from '@angular/fire/auth';
import { Observable, BehaviorSubject, firstValueFrom, from } from 'rxjs';
import { filter, map, switchMap } from 'rxjs/operators';
import {Firestore, doc, getDoc, collection, getDocs, setDoc} from '@angular/fire/firestore';
import { HttpClient } from '@angular/common/http';
import { environment } from '../environment/environment';
import { UserData } from '../Models/UserData';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private currentUserSubject = new BehaviorSubject<User | null>(null);
  private userDataSubject = new BehaviorSubject<UserData | null>(null);
  private idTenantSubject = new BehaviorSubject<string | null>(null);
  url:string = 'http://localhost:80/'
  constructor(private afAuth: Auth, private firestore: Firestore, private http: HttpClient, private injector: Injector) {
    authState(this.afAuth).subscribe(async (user) => {
      this.currentUserSubject.next(user);
      if (user) {
          const fallbackUserData: UserData = {
            id: '', // Se llenará desde el backend
            firebaseUID: user.uid,
            email: user.email || '',
            name: user.displayName || user.email?.split('@')[0] || '',
            rol: 'cliente', // rol por defecto
            profile: {
              providerId: '',
              uid: user.uid,
              displayName: user.displayName,
              email: user.email,
              phoneNumber: user.phoneNumber,
              photoURL: user.photoURL
            }
          };
          this.userDataSubject.next(fallbackUserData);
      } else {
        this.userDataSubject.next(null);
      }
    });
  }
  /**
   * Maneja el inicio de sesión del usuario con correo y contraseña.
   * @param email correo electrónico del usuario
   * @param password contraseña del usuario 
   */
  async login(email: string, password: string, idTenant:string) {
    await runInInjectionContext(this.injector, async () => {
      await signInWithEmailAndPassword(this.afAuth, email, password);
      try{
        await this.fetchCurrentUserFromBackend(idTenant);
      }catch(error){
        console.error('Error al cargar datos del usuario tras lgon: ', error);
      }
    });
  }
  /**
   * Consulta al backend los datos del usuario actualmente autenticado y los guarda en userData.
   * @returns Los datos del usuario actual obtenidos desde el backend
   */
  async fetchCurrentUserFromBackend(idTenant: string) {
    const token = await this.getToken();
    this.idTenantSubject.next(idTenant);
    console.log('Token obtenido de Firebase:', token); 
  
    if (token) {
      const headers = { 
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      };
      
      try {
        const backendUserData = await firstValueFrom(
          this.http.get<UserData>(this.url + idTenant + '/api/users/me', { headers })
        );
        
        // Guardar los datos del backend en userDataSubject
        this.userDataSubject.next(backendUserData);
        console.log('Datos del usuario desde backend cargados y guardados:', backendUserData);
        
        return backendUserData;
      } catch (error) {
        console.error('Error al consultar el backend:', error);
        throw error;
      }
    } else {
      console.warn('No se obtuvo token, no se puede consultar el backend');
      throw new Error('No se pudo obtener el token de autenticación');
    }
  }
  /**
   * Observable que emite true si el usuario está autenticado, false en caso contrario.
   */
  get isLoggedIn$(): Observable<boolean> {
    return this.currentUserSubject.pipe(map(user => !!user));
  }
  /**
   * Observable que trae el usuario actualmente autenticado, o null si no hay ninguno.
   */
  get currentUser(): Observable<User | null> {
    return this.currentUserSubject.asObservable();
  }
  /**
   * Observable que emite el UID del usuario autenticado, o null si no hay ninguno.
   */
  get uid(): Observable<string | null> {
    return this.currentUser.pipe(map(user => user?.uid ?? null));
  }
  /**
   *  Obtiene los claims personalizados del usuario autenticado.
   * @returns Un observable que emite los claims personalizados del usuario autenticado, o null si no hay ninguno.
   */
  getUserClaims(): Observable<{ [key: string]: any } | null> {
    return this.currentUser.pipe(
      filter((user): user is User => !!user),
      switchMap(user => from(getIdTokenResult(user, true))),
      map(idTokenResult => idTokenResult.claims)
    );
  }
  /**
   * Obtiene el token de ID del usuario autenticado.
   * @returns Una promesa que resuelve con el token de ID, o null si no hay usuario autenticado.
   */
  async getToken(): Promise<string | null> {
    const immediateUser = await this.afAuth.currentUser;
    if (immediateUser) {
      return immediateUser.getIdToken(true);
    }
    const user = await firstValueFrom(
      this.currentUser.pipe(
        filter((u): u is User => u !== null)
      )
    );
    return user.getIdToken(true);
  }
  /**
   * 
   * @returns Un observable que emite true si el usuario tiene el rol de admin, false en caso contrario.
   */
  isAdmin(): Observable<boolean> {
        return this.getUserClaims().pipe(
          map(claims => !!(claims && claims['admin'] === true))
        );
  }
  /*
  * Cierra la sesión del usuario actualmente autenticado.
  */
  logout() {
    return signOut(this.afAuth);
  }
  
  get idTenant$(): Observable<string | null> {
    return this.idTenantSubject.asObservable();
  }
  /**
   * Observable que trae los datos completos del usuario desde el backend
   */
  get userData(): Observable<UserData | null> {
    return this.userDataSubject.asObservable();
  }
}