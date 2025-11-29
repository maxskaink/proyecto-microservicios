import { Injectable, Injector, runInInjectionContext } from '@angular/core';
import {
  Auth,
  signInWithEmailAndPassword,
  signInWithPopup,
  GoogleAuthProvider,
  signOut,
  User,
  onAuthStateChanged,
  getIdTokenResult,
  authState,
  getIdToken
} from '@angular/fire/auth';
import { BehaviorSubject, Observable, from, firstValueFrom, combineLatest } from 'rxjs';
import { map, filter, switchMap } from 'rxjs/operators';
import { Firestore } from '@angular/fire/firestore';
import { HttpClient } from '@angular/common/http';

import { UserData } from '../Models/UserData';
import { TenantService } from './TenantService';


@Injectable({ providedIn: 'root' })
export class AuthService {

  private currentUserSubject = new BehaviorSubject<User | null>(null);
  private userDataSubject = new BehaviorSubject<UserData | null>(null);
  public userCurrentData: UserData | null = null;
  private authReadySubject = new BehaviorSubject<boolean>(false);
  authReady$ = this.authReadySubject.asObservable();

  url: string = 'http://localhost:80/';

  private cachedToken: string | null = null;
  private tokenExpiry: number | null = null;
  private tokenPromise: Promise<string | null> | null = null;

  constructor(
    private afAuth: Auth,
    private firestore: Firestore,
    private http: HttpClient,
    private injector: Injector,
    private tenantService: TenantService
  ) {
      onAuthStateChanged(this.afAuth, async (user) => {
        console.log('🔄 Auth state changed:', user?.email || 'no user');

        this.currentUserSubject.next(user);

        
        this.cachedToken = null;
        this.tokenExpiry = null;
        this.tokenPromise = null;

        
        this.authReadySubject.next(true);

        if (user) {
          const savedTenant = this.tenantService.getTenant();
          if (savedTenant) {
            this.tenantService.setTenant(savedTenant);
            await this.loadUserData(savedTenant);
          }
        } else {
          this.userDataSubject.next(null);
          this.tenantService.clearTenant();
        }
      });
  }

  async login(email: string, password: string, idTenant: string) {
    await runInInjectionContext(this.injector, async () => {
      await signInWithEmailAndPassword(this.afAuth, email, password);

      this.tenantService.setTenant(idTenant);

      await this.loadUserData(idTenant);
    });
  }

  async loginWithGoogle(idTenant: string) {
    await runInInjectionContext(this.injector, async () => {
      const provider = new GoogleAuthProvider();
      
      // Configurar el provider para forzar la selección de cuenta
      provider.setCustomParameters({
        prompt: 'select_account'
      });

      // Realizar login con popup
      const result = await signInWithPopup(this.afAuth, provider);
      
      if (result.user) {
        this.tenantService.setTenant(idTenant);
        await this.loadUserData(idTenant);
        
        console.log('✅ Login con Google exitoso:', result.user.email);
        return result.user;
      }
      
      throw new Error('No se pudo completar el login con Google');
    });
  }

  private async loadUserData(idTenant: string): Promise<void> {
    try {
      const token = await this.getToken();
      if (!token) return;

      const backendUserData = await firstValueFrom(
        this.http.get<UserData>(`${this.url}${idTenant}/api/users/me`, {
          headers: { Authorization: `Bearer ${token}` }
        })
      );
      this.userCurrentData = backendUserData;
      this.userDataSubject.next(backendUserData);
      localStorage.setItem('user_data', JSON.stringify(backendUserData));
      console.log('✅ userData cargado:', backendUserData);

    } catch (error) {
      console.error('❌ Error al cargar userData:', error);
    }
  }


  async getToken(): Promise<string | null> {
    // Espera a que Firebase emita el usuario autenticado
    const user = await firstValueFrom(authState(this.afAuth));

    if (!user) return null;

    // Obtener token usando función modular
    const token = await getIdToken(user);

    // Opcional: guardar token para acelerar carga
    localStorage.setItem('auth_token', token);

    return token;
  }

  get isLoggedIn$(): Observable<boolean> {
    return combineLatest([
      this.currentUserSubject,
      this.authReady$
    ]).pipe(
      filter(([_, ready]) => ready),   // 👈 esperar a que Firebase responda
      map(([user]) => !!user)
    );
  }


  get currentUser(): Observable<User | null> {
    return this.currentUserSubject.asObservable();
  }
  getCurrentUser(): BehaviorSubject<User | null> {
    return this.currentUserSubject;
  }
  get userData(): Observable<UserData | null> {
    return this.userDataSubject.asObservable();
  }

  async logout() {
    this.tenantService.clearTenant();
    this.cachedToken = null;
    this.tokenExpiry = null;
    this.tokenPromise = null;

    return signOut(this.afAuth);
  }

  getUserClaims() {
    return this.currentUser.pipe(
      filter((u): u is User => !!u),
      switchMap(user => from(getIdTokenResult(user, true))),
      map(result => result.claims)
    );
  }
   getUserRole(): string | null {
    return this.userCurrentData?.rol ?? null;
  }

  isAdmin() {
    return this.getUserClaims().pipe(
      map(claims => claims?.['admin'] === true)
    );
  }
}
