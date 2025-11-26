import { Injectable, Injector, runInInjectionContext } from '@angular/core';
import {
  Auth,
  signInWithEmailAndPassword,
  signOut,
  User,
  onAuthStateChanged,
  getIdTokenResult
} from '@angular/fire/auth';
import { BehaviorSubject, Observable, from, firstValueFrom } from 'rxjs';
import { map, filter, switchMap } from 'rxjs/operators';
import { Firestore } from '@angular/fire/firestore';
import { HttpClient } from '@angular/common/http';

import { UserData } from '../Models/UserData';
import { TenantService } from './TenantService';

@Injectable({ providedIn: 'root' })
export class AuthService {

  private currentUserSubject = new BehaviorSubject<User | null>(null);
  private userDataSubject = new BehaviorSubject<UserData | null>(null);

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

  private async loadUserData(idTenant: string): Promise<void> {
    try {
      const token = await this.getToken();
      if (!token) return;

      const backendUserData = await firstValueFrom(
        this.http.get<UserData>(`${this.url}${idTenant}/api/users/me`, {
          headers: { Authorization: `Bearer ${token}` }
        })
      );

      this.userDataSubject.next(backendUserData);
      console.log('✅ userData cargado:', backendUserData);

    } catch (error) {
      console.error('❌ Error al cargar userData:', error);
    }
  }

  async getToken(): Promise<string | null> {
    const user = this.afAuth.currentUser;
    if (!user) return null;

    const now = Date.now();

    if (this.cachedToken && this.tokenExpiry && now < this.tokenExpiry) {
      return this.cachedToken;
    }

    if (this.tokenPromise) {
      return this.tokenPromise;
    }

    this.tokenPromise = user.getIdTokenResult(true)
      .then(res => {
        this.cachedToken = res.token;
        this.tokenExpiry = new Date(res.expirationTime).getTime();
        return this.cachedToken;
      })
      .finally(() => this.tokenPromise = null);

    return this.tokenPromise;
  }

  get isLoggedIn$(): Observable<boolean> {
    return this.currentUserSubject.pipe(map(user => !!user));
  }

  get currentUser(): Observable<User | null> {
    return this.currentUserSubject.asObservable();
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

  isAdmin() {
    return this.getUserClaims().pipe(
      map(claims => claims?.['admin'] === true)
    );
  }
}
