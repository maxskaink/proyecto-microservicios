import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { AuthService } from '../../../service /Authser.vice';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-header',
  imports: [CommonModule],
  templateUrl: './header.html',
  styleUrl: './header.css',
})
export class Header implements OnInit, OnDestroy {
  
  // Estado del menú móvil
  isMenuOpen: boolean = false;
  
  // Estado de autenticación
  isLoggedIn: boolean = false;
  currentUser: any = null;
  
  // Contador de carrito (simulado por ahora)
  cartItemCount: number = 0;
  
  // Suscripciones
  private subscriptions = new Subscription();

  constructor(private router: Router, private authService: AuthService) {}

  ngOnInit(): void {
    this.initializeAuth();
  }

  ngOnDestroy(): void {
    this.subscriptions.unsubscribe();
  }

  // Inicializar autenticación
  private initializeAuth(): void {
    // Suscribirse al estado de autenticación
    const authSub = this.authService.isLoggedIn$.subscribe(loggedIn => {
      this.isLoggedIn = loggedIn;
    });
    this.subscriptions.add(authSub);

    // Suscribirse a los datos del usuario
    const userSub = this.authService.userData.subscribe(userData => {
      this.currentUser = userData;
    });
    this.subscriptions.add(userSub);
  }

  // Toggle del menú móvil
  toggleMenu(): void {
    this.navigateTo('/user')
  }

  // Cerrar menú móvil
  closeMenu(): void {
    this.isMenuOpen = false;
  }

  // Navegación
  navigateTo(route: string): void {
    this.router.navigate([route]);
    this.closeMenu();
  }

  // Ir al perfil de usuario
  goToUserProfile(): void {
    if (this.isLoggedIn) {
      this.navigateTo('/user');
    } else {
      this.navigateTo('/login');
    }
  }

  // Ir al carrito
  goToCart(): void {
    this.navigateTo('/cart');
  }

  // Cerrar sesión
  async logout(): Promise<void> {
    try {
      await this.authService.logout();
      this.navigateTo('/login');
      this.closeMenu();
    } catch (error) {
      console.error('Error al cerrar sesión:', error);
    }
  }

  // Obtener nombre del usuario
  getUserName(): string {
    if (!this.currentUser) return 'Usuario';
    
    return this.currentUser.nombre || 
           this.currentUser.displayName ||
           this.currentUser.email?.split('@')[0] || 
           'Usuario';
  }

  // Simular actualización del carrito (conectar con servicio real después)
  updateCartCount(count: number): void {
    this.cartItemCount = count;
  }

  // Método para testing - agregar item al carrito
  addToCart(): void {
    this.cartItemCount++;
  }

  // Método para testing - limpiar carrito
  clearCart(): void {
    this.cartItemCount = 0;
  }
}
