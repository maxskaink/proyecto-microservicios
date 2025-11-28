import { Injectable } from '@angular/core';
import { CanActivate, Router, UrlTree } from '@angular/router';
import { Observable, combineLatest, of } from 'rxjs';
import { filter, map, take, catchError } from 'rxjs/operators';
import { AuthService } from '../service/Authser.vice';
import { TenantService } from '../service/TenantService';

@Injectable({ providedIn: 'root' })
export class authGuard implements CanActivate {

  constructor(
    private authService: AuthService,
    private tenantService: TenantService,
    private router: Router
  ) {}

  canActivate(): Observable<boolean | UrlTree> {

    return combineLatest([
      this.authService.authReady$,          // 🚀 Espera a que Firebase termine
      this.authService.currentUser,         // Usuario real
      this.tenantService.tenant$            // Tenant
    ]).pipe(
      // SOLO avanzar cuando authReady$ === true
      filter(([ready]) => ready === true),

      take(1),

      map(([_, user, tenantId]) => {

        if (!user) {
          console.log("❌ Usuario NO autenticado → login");
          return this.router.createUrlTree(['/login']);
        }

        if (!tenantId) {
          console.log("⚠ Usuario autenticado pero SIN tenant → register-tenant");
          return this.router.createUrlTree(['/register-tenant']);
        }

        console.log("✔ Usuario y tenant OK:", tenantId);
        return true;
      }),

      catchError((err) => {
        console.error("❌ Error en authGuard:", err);
        return of(this.router.createUrlTree(['/login']));
      })
    );
  }
}
