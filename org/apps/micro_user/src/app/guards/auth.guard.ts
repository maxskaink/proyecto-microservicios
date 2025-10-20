import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { SimpleAuthService } from '@org/shared';

@Injectable({
  providedIn: 'root'
})
export class AuthGuard {
  
  constructor(
    private authService: SimpleAuthService,
    private router: Router
  ) {}

  canActivate(): boolean {
    if (this.authService.isAuthenticated()) {
      return true;
    } else {
      this.router.navigate(['/login']);
      return false;
    }
  }
}
