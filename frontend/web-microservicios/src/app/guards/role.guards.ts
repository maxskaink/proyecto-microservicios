import { Injectable } from '@angular/core';
import { CanActivate, Router, UrlTree } from '@angular/router';
import { Observable, combineLatest, of } from 'rxjs';
import { filter, map, take, catchError } from 'rxjs/operators';
import { AuthService } from '../service/Authser.vice';
import { TenantService } from '../service/TenantService';

@Injectable({ providedIn: 'root' })
export class roleGuard implements CanActivate {

  constructor(
    private authService: AuthService, 
    private tenantService: TenantService,
    private router: Router
  ) {}

  canActivate(): Observable<boolean | UrlTree> {
    return combineLatest([
      this.authService.isLoggedIn$,
      this.tenantService.tenant$
    ]).pipe(
      take(1),
      map(([isLoggedIn, tenantId]) => {
        
        // Verificar si está autenticado
        if (!isLoggedIn) {
          console.log('❌ Usuario no autenticado, redirigiendo a login');
          return this.router.createUrlTree(['/login']);
        }

        // Verificar si tiene tenant seleccionado
        if (!tenantId) {
          console.log('⚠️ Usuario sin tenant, redirigiendo a selección de tenant');
          return this.router.createUrlTree(['/register-tenant']);
        }

        console.log('✅ Usuario autenticado y con tenant:', tenantId);
        return true;
      }),
      catchError((error) => {
        console.error('❌ Error en roleGuard:', error);
        return of(this.router.createUrlTree(['/login']));
      })
    );
  }
}
