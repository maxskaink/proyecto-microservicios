import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { AuthService } from '../service/Authser.vice';
import { from, mergeMap } from 'rxjs';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const auth = inject(AuthService);

  // Convertimos la promesa del token a Observable
  return from(auth.getToken()).pipe(
    mergeMap((token) => {
      const newReq = token
        ? req.clone({ setHeaders: { Authorization: `Bearer ${token}` } })
        : req;

      return next(newReq);
    })
  );
};
