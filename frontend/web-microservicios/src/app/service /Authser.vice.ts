import { Injectable, Injector, runInInjectionContext } from '@angular/core';
import { Auth, signInWithEmailAndPassword, signOut, User, authState, getIdTokenResult } from '@angular/fire/auth';
import { Observable, BehaviorSubject, firstValueFrom, from } from 'rxjs';
import { filter, map, switchMap } from 'rxjs/operators';
import {Firestore, doc, getDoc, collection, getDocs, setDoc} from '@angular/fire/firestore';
import { HttpClient } from '@angular/common/http';
import { environment } from '../environment/environment';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private currentUserSubject = new BehaviorSubject<User | null>(null);
  private userDataSubject = new BehaviorSubject<any | null>(null);

  constructor(private afAuth: Auth, private firestore: Firestore, private http: HttpClient, private injector: Injector) {
    authState(this.afAuth).subscribe(async (user) => {
      this.currentUserSubject.next(user);
      if (user) {
        try {
          const docSnap = await runInInjectionContext(this.injector, () => 
            getDoc(doc(this.firestore, 'usuarios', user.uid))
          );
          if (docSnap.exists()) {
            const firestoreData = docSnap.data();
            // Combinar datos de Firebase Auth con Firestore
            const combinedData = {
              uid: user.uid,
              email: user.email,
              emailVerified: user.emailVerified,
              displayName: user.displayName,
              ...firestoreData
            };
            this.userDataSubject.next(combinedData);
            console.log('Datos completos del usuario:', combinedData);
          } else {
            console.warn('Usuario autenticado pero sin documento en Firestore');
            this.userDataSubject.next({
              uid: user.uid,
              email: user.email,
              emailVerified: user.emailVerified,
              displayName: user.displayName
            });
          }
        } catch (error) {
          console.error('Error al cargar datos de Firestore:', error);
          this.userDataSubject.next({
            uid: user.uid,
            email: user.email,
            emailVerified: user.emailVerified,
            displayName: user.displayName
          });
        }
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
  async login(email: string, password: string) {
    await runInInjectionContext(this.injector, async () => {
      await signInWithEmailAndPassword(this.afAuth, email, password);
    });
  }
  /**
   * Consulta al backend los datos del usuario actualmente autenticado.
   * @returns Los datos del usuario actual obtenidos desde el backend
   */
  async fetchCurrentUserFromBackend() {
    const token = await this.getToken();
    console.log('Token obtenido de Firebase:', token); 
  
    const url = `${environment.apiUrl}/auth/me`;
    if (token) {
      return firstValueFrom(
        this.http.get(url, {
          headers: { Authorization: `Bearer ${token}` },
        })
      );
    }
    console.warn('No se obtuvo token, enviando sin Authorization header');
    return firstValueFrom(this.http.get(url));
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

  /**
   * Obtiene los datos completos del usuario desde Firestore
   * @returns Los datos del usuario almacenados en Firestore
   */
  async getUserDataFromFirestore(): Promise<any | null> {
    const user = await firstValueFrom(
      this.currentUser.pipe(
        filter((u): u is User => u !== null)
      )
    );
    
    if (!user) return null;

    try {
      const docSnap = await runInInjectionContext(this.injector, () => 
        getDoc(doc(this.firestore, 'usuarios', user.uid))
      );
      
      if (docSnap.exists()) {
        return docSnap.data();
      } else {
        console.warn('No se encontró documento del usuario en Firestore');
        return null;
      }
    } catch (error) {
      console.error('Error al obtener datos de Firestore:', error);
      return null;
    }
  }

  /**
   * Observable que combina datos de Firebase Auth con datos de Firestore
   * @returns Observable con datos completos del usuario
   */
  get fullUserData(): Observable<any | null> {
    return this.currentUser.pipe(
      switchMap(async (user) => {
        if (!user) return null;
        
        try {
          // Datos básicos de Firebase Auth
          const authData = {
            uid: user.uid,
            email: user.email,
            emailVerified: user.emailVerified,
            displayName: user.displayName
          };

          // Datos adicionales de Firestore
          const firestoreData = await runInInjectionContext(this.injector, async () => {
            const docSnap = await getDoc(doc(this.firestore, 'usuarios', user.uid));
            return docSnap.exists() ? docSnap.data() : {};
          });

          // Combinar ambos
          return {
            ...authData,
            ...firestoreData
          };
        } catch (error) {
          console.error('Error al combinar datos del usuario:', error);
          return {
            uid: user.uid,
            email: user.email,
            emailVerified: user.emailVerified,
            displayName: user.displayName
          };
        }
      })
    );
  }

  /**
   * Observable que trae los datos completos del usuario (Firebase Auth + Firestore)
   */
  get userData(): Observable<any | null> {
    return this.userDataSubject.asObservable();
  }
}