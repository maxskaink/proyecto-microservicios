import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { SimpleAuthService } from '@org/shared';
import { User } from '@angular/fire/auth';

@Component({
  selector: 'app-dashboard',
  imports: [CommonModule],
  template: `
    <div class="container mt-5">
      <div class="row justify-content-center">
        <div class="col-md-8">
          <div class="card">
            <div class="card-header">
              <h3>Dashboard - ¡Bienvenido!</h3>
            </div>
            <div class="card-body">
              <div *ngIf="user" class="mb-3">
                <h5>Información del Usuario:</h5>
                <p><strong>Email:</strong> {{ user.email }}</p>
                <p><strong>UID:</strong> {{ user.uid }}</p>
                <p><strong>Último login:</strong> {{ user.metadata.lastSignInTime | date:'medium' }}</p>
              </div>
              
              <div class="d-grid gap-2">
                <button class="btn btn-danger" (click)="logout()">
                  Cerrar Sesión
                </button>
                <button class="btn btn-secondary" (click)="goToLogin()">
                  Volver al Login
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .card {
      box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
      border-radius: 10px;
    }
    .card-header {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
      border-radius: 10px 10px 0 0;
    }
  `]
})
export class Dashboard implements OnInit {
  user: User | null = null;

  constructor(
    private authService: SimpleAuthService,
    private router: Router
  ) {}

  ngOnInit() {
    // Obtener usuario actual
    this.user = this.authService.getCurrentUser();
    
    // Suscribirse a cambios de autenticación
    this.authService.currentUser$.subscribe(user => {
      this.user = user;
      if (!user) {
        this.router.navigate(['/login']);
      }
    });
  }

  async logout() {
    try {
      await this.authService.logout();
      console.log('Logout exitoso');
      this.router.navigate(['/login']);
    } catch (error) {
      console.error('Error en logout:', error);
    }
  }

  goToLogin() {
    this.router.navigate(['/login']);
  }
}
