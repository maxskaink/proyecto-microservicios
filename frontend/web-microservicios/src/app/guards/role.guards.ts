import { Injectable } from '@angular/core';
import { CanActivate, ActivatedRouteSnapshot, Router, UrlTree } from '@angular/router';
import { AuthService } from '../service/Authser.vice';


@Injectable({
  providedIn: 'root'
})
export class RoleGuard implements CanActivate {

  constructor(
    private authService: AuthService,
    private router: Router
  ) {}

  canActivate(route: ActivatedRouteSnapshot): boolean | UrlTree {
    const allowedRoles: string[] = route.data['roles'];  // roles permitidos en la ruta
    const userRole = this.authService.userCurrentData?.rol; // rol del usuario (ajusta según tu modelo)

    if (!userRole) {
      console.warn('⛔ No hay usuario autenticado');
      return this.router.parseUrl('/login');
    }

    if (allowedRoles.includes(userRole)) {
      return true; // tiene permiso
    }

    console.warn(`⛔ Acceso denegado: se requiere uno de estos roles: ${allowedRoles}`);
    return this.router.parseUrl('/forbidden'); // redirige a página de acceso denegado
  }
}
