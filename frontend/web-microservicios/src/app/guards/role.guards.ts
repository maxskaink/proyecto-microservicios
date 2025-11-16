import { Injectable } from '@angular/core';
import { CanActivate, Router, UrlTree } from '@angular/router';
import { Observable } from 'rxjs';

import { filter, map, switchMap, take } from 'rxjs/operators';
import { AuthService } from '../service/Authser.vice';

@Injectable({ providedIn: 'root' })
export class roleGuard implements CanActivate {

  constructor(private auth: AuthService, private router: Router) {}

  canActivate(): Observable<boolean | UrlTree> {
    
    return this.auth.authReady$.pipe(

      // ⏳ Espera a Firebase
      filter(ready => ready === true),
      take(1),

      // ✔ Ahora sí revisa si está logeado
      switchMap(() => this.auth.isLoggedIn$),

      map(isLogged => {
        if (isLogged) return true;
        return this.router.createUrlTree(['/login']);
      })
    );
  }
}
