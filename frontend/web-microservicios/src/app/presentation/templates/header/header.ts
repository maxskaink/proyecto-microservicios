import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService } from '../../../service/Authser.vice';
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
  
  
  
  // Contador de carrito
  cartItemCount: number = 0;
  
  // Suscripciones
  private subscriptions = new Subscription();

  constructor(
    private router: Router, 
    private authService: AuthService,
  ) {}

  ngOnInit(): void {
    this.initializeAuth();
    this.getUserName();
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



  // Obtener nombre del usuario
  getUserName(): string {
    if (!this.currentUser) return 'Usuario';
    
    return this.currentUser.name || 
           this.currentUser.email?.split('@')[0] || 
           'Usuario';
  }


}
