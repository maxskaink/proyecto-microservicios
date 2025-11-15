import { Injectable, Injector, runInInjectionContext } from '@angular/core';
import { 
  Auth, 
  signInWithEmailAndPassword, 
  signOut, 
  User, 
  onAuthStateChanged, 
  getIdTokenResult 
} from '@angular/fire/auth';
import { BehaviorSubject, Observable, firstValueFrom, from } from 'rxjs';
import { filter, map, switchMap } from 'rxjs/operators';
import { Firestore } from '@angular/fire/firestore';
import { HttpClient } from '@angular/common/http';
import { UserData } from '../Models/UserData';

@Injectable({ providedIn: 'root' })
export class AuthService {

  private currentUserSubject = new BehaviorSubject<User | null>(null);
  private userDataSubject = new BehaviorSubject<UserData | null>(null);
  private idTenantSubject = new BehaviorSubject<string | null>(null);

  /** ⛔ NUEVO: indica cuando Firebase terminó de restaurar el usuario */
  private authReadySubject = new BehaviorSubject<boolean>(false);
  authReady$ = this.authReadySubject.asObservable();

  url: string = 'http://localhost:80/';

  constructor(
    private afAuth: Auth,
    private firestore: Firestore,
    private http: HttpClient,
    private injector: Injector
  ) {

    // 🔥 Restauración automática DEBE completarse ANTES que el guard
    onAuthStateChanged(this.afAuth, async (user) => {
      console.log('🔄 Auth state changed:', user?.email || 'no user');

      this.currentUserSubject.next(user);

      if (user) {
        const savedTenant = this.getSavedTenant();
        if (savedTenant) {
          console.log('📦 Tenant guardado encontrado:', savedTenant);
          await this.loadUserData(savedTenant);
        }
      } else {
        console.log('❌ No hay usuario, limpiando datos');
        this.userDataSubject.next(null);
        this.idTenantSubject.next(null);
      }

      // 🔥 Avisar a los guards que Firebase YA terminó
      this.authReadySubject.next(true);
    });
  }

  async login(email: string, password: string, idTenant: string) {
    await runInInjectionContext(this.injector, async () => {
      await signInWithEmailAndPassword(this.afAuth, email, password);
      this.saveTenant(idTenant);
      await this.loadUserData(idTenant);
    });
  }

  private async loadUserData(idTenant: string): Promise<void> {
    try {
      const token = await this.afAuth.currentUser?.getIdToken(true);
      if (!token) return;

      this.idTenantSubject.next(idTenant);

      const backendUserData = await firstValueFrom(
        this.http.get<UserData>(`${this.url}${idTenant}/api/users/me`, {
          headers: { 'Authorization': `Bearer ${token}` }
        })
      );

      this.userDataSubject.next(backendUserData);
      console.log('✅ userData cargado:', backendUserData);

    } catch (error) {
      console.error('❌ Error al cargar userData:', error);
    }
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

  get idTenant$(): Observable<string | null> {
    return this.idTenantSubject.asObservable();
  }

  async getToken(): Promise<string | null> {
    const user = await this.afAuth.currentUser;
    return user ? await user.getIdToken(true) : null;
  }

  private saveTenant(tenantId: string): void {
    localStorage.setItem('currentTenant', tenantId);
  }

  private getSavedTenant(): string | null {
    return localStorage.getItem('currentTenant');
  }

  async logout() {
    localStorage.removeItem('currentTenant');
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
