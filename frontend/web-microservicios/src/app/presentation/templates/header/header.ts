import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService } from '../../../service /Authser.vice';
import { ShoppingCartService } from '../../../service /ShoppinCartService';
import { Subscription } from 'rxjs';
import { UserData } from '../../../Models/UserData';

@Component({
  selector: 'app-header',
  imports: [CommonModule, FormsModule],
  templateUrl: './header.html',
  styleUrl: './header.css',
})
export class Header implements OnInit, OnDestroy {
  
  isLoggedIn: boolean = false;
  currentUser: UserData | null = null;
  
  // Menú móvil
  isMenuOpen: boolean = false;
  
  // Búsqueda
  searchTerm: string = '';
  
  
  // Contador de carrito
  cartItemCount: number = 0;
  
  // Suscripciones
  private subscriptions = new Subscription();

  constructor(
    private router: Router, 
    private authService: AuthService,
    private shoppingCartService: ShoppingCartService
  ) {}

  ngOnInit(): void {
    this.initializeAuth();
    this.getUserName();
    this.loadCartItemCount();
  }

  ngOnDestroy(): void {
    this.subscriptions.unsubscribe();
  }

  // Inicializar autenticación
  private initializeAuth(): void {
    // Suscribirse al estado de autenticación
    const authSub = this.authService.isLoggedIn$.subscribe((loggedIn: boolean) => {
      this.isLoggedIn = loggedIn;
    });
    this.subscriptions.add(authSub);

    // Suscribirse a los datos del usuario
    const userSub = this.authService.userData.subscribe((userData: UserData | null) => {
      this.currentUser = userData;
    });
    this.subscriptions.add(userSub);
  }

  // Toggle del menú móvil
  toggleMenu(): void {
    this.isMenuOpen = !this.isMenuOpen;
    // También navegar al perfil de usuario
    this.navigateTo('/user');
  }

  // Navegación
  navigateTo(route: string): void {
    this.router.navigate([route]);
  }

  // Ir al perfil de usuario
  goToUserProfile(): void {
    if (this.isLoggedIn) {
      this.navigateTo('/user');
    } else {
      this.navigateTo('/login');
    }
  }
  goToHome(): void {
    this.navigateTo('/home');
  }
  // Ir al carrito
  goToCart(): void {
    this.navigateTo('/shopping-cart');
  }

  // Función de búsqueda
  onSearch(): void {
    if (this.searchTerm.trim()) {
      console.log('Buscando:', this.searchTerm);
      // Navegar a home con parámetro de búsqueda
      this.router.navigate(['/home'], { 
        queryParams: { search: this.searchTerm.trim() } 
      });
    }
  }

  // Obtener nombre del usuario
  getUserName(): string {
    if (!this.currentUser) return 'Usuario';
    
    return this.currentUser.name || 
           this.currentUser.email?.split('@')[0] || 
           'Usuario';
  }

  // Cargar contador del carrito
  loadCartItemCount(): void {
    if (this.isLoggedIn) {
      const cartSub = this.shoppingCartService.getCartItemCount().subscribe({
        next: (count) => {
          this.cartItemCount = count;
        },
        error: (error) => {
          console.error('Error al cargar contador del carrito:', error);
          this.cartItemCount = 0;
        }
      });
      this.subscriptions.add(cartSub);
    }
  }

  // Actualizar contador del carrito
  updateCartCount(count: number): void {
    this.cartItemCount = count;
  }
}
