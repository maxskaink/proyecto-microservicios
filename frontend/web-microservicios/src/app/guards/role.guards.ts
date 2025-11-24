import { Injectable } from '@angular/core';
import { CanActivate, Router, UrlTree } from '@angular/router';
import { Observable, combineLatest } from 'rxjs';
import { filter, map, take } from 'rxjs/operators';
import { AuthService } from '../service/Authser.vice';

@Injectable({ providedIn: 'root' })
export class roleGuard implements CanActivate {

  constructor(private auth: AuthService, private router: Router) {}

  canActivate(): Observable<boolean | UrlTree> {
    return combineLatest([
      this.auth.authReady$,
      this.auth.isLoggedIn$
    ]).pipe(
      // Espera a que Firebase termine de restaurar la sesión
      filter(([ready, logged]) => ready),
      take(1), // solo necesitamos la primera emisión después de authReady
      map(([_, logged]) => {
        if (logged) return true;
        return this.router.createUrlTree(['/login']);
      })
    );
  }
}
