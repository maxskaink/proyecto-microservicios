import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { AuthService } from '../../../service /Authser.vice';
import { Subscription } from 'rxjs';

interface MenuOption {
  id: string;
  title: string;
  description: string;
  icon: string;
  route: string;
  category: 'shopping' | 'profile' | 'business';
}

interface UserData {
  // Campos de Firebase Auth
  uid?: string;
  email?: string;
  emailVerified?: boolean;
  displayName?: string;
  
  // Campos personalizados de Firestore
  nombre?: string;
  apellido?: string;
  telefono?: string;
  direccion?: string;
  rol?: string;
  fechaRegistro?: string;
  activo?: boolean;
}

@Component({
  selector: 'app-user',
  imports: [CommonModule],
  templateUrl: './user.html',
  styleUrl: './user.css',
})
export class User implements OnInit, OnDestroy {
  
  // Información del usuario actual
  currentUser: UserData | null = null;
  isLoggedIn: boolean = false;
  isLoading: boolean = true;
  
  // Suscripciones para limpiar al destruir el componente
  private subscriptions: Subscription = new Subscription();
  
  menuOptions: MenuOption[] = [
    // Sección de Compras
    {
      id: 'catalog',
      title: 'Catálogo',
      description: 'Explora todos los productos disponibles',
      icon: 'bi-grid-3x3-gap',
      route: '/catalog',
      category: 'shopping'
    },
    {
      id: 'cart',
      title: 'Carrito de Compras',
      description: 'Ver productos en tu carrito',
      icon: 'bi-cart3',
      route: '/cart',
      category: 'shopping'
    },
    {
      id: 'purchase-history',
      title: 'Historial de Compras',
      description: 'Revisa tus compras anteriores',
      icon: 'bi-clock-history',
      route: '/purchase-history',
      category: 'shopping'
    },
    // Sección de Perfil
    {
      id: 'edit-profile',
      title: 'Editar Perfil',
      description: 'Actualiza tu información personal',
      icon: 'bi-person-gear',
      route: '/edit-profile',
      category: 'profile'
    },
    // Sección de Negocio
    {
      id: 'my-products',
      title: 'Mis Productos',
      description: 'Gestiona los productos que vendes',
      icon: 'bi-box-seam',
      route: '/my-products',
      category: 'business'
    },
    {
      id: 'publish-product',
      title: 'Publicar Producto',
      description: 'Añade un nuevo producto a la venta',
      icon: 'bi-plus-circle',
      route: '/publish-product',
      category: 'business'
    },
    {
      id: 'orders',
      title: 'Pedidos',
      description: 'Gestiona los pedidos de tus productos',
      icon: 'bi-clipboard-check',
      route: '/orders',
      category: 'business'
    }
  ];

  constructor(private router: Router, private authService: AuthService) {}

  ngOnInit(): void {
    this.loadUserData();
  }

  ngOnDestroy(): void {
    this.subscriptions.unsubscribe();
  }

  // Cargar datos del usuario autenticado
  private loadUserData(): void {
    this.isLoading = true;

    // Suscribirse al estado de autenticación
    const authSub = this.authService.isLoggedIn$.subscribe(loggedIn => {
      this.isLoggedIn = loggedIn;
      if (!loggedIn) {
        this.currentUser = null;
        this.isLoading = false;
      }
    });
    this.subscriptions.add(authSub);

    // Obtener datos completos del usuario (Firebase Auth + Firestore)
    const userDataSub = this.authService.userData.subscribe(userData => {
      if (userData) {
        this.currentUser = userData as UserData;
        console.log('Datos del usuario cargados:', this.currentUser);
        this.isLoading = false;
      } else if (!this.isLoggedIn) {
        this.currentUser = null;
        this.isLoading = false;
      }
    });
    this.subscriptions.add(userDataSub);

    // Obtener claims para verificar roles (como fallback si no está en Firestore)
    const claimsSub = this.authService.getUserClaims().subscribe(claims => {
      if (claims && this.currentUser && !this.currentUser.rol) {
        this.currentUser.rol = claims['admin'] ? 'admin' : 'usuario';
      }
    });
    this.subscriptions.add(claimsSub);

    // Si después de 5 segundos sigue cargando, detener el loading
    setTimeout(() => {
      if (this.isLoading) {
        this.isLoading = false;
      }
    }, 5000);
  }

  // Recargar datos del usuario desde Firestore
  async loadAdditionalUserData(): Promise<void> {
    if (!this.isLoggedIn) return;
    
    try {
      console.log('Recargando datos del usuario desde Firestore...');
      const userData = await this.authService.getUserDataFromFirestore();
      if (userData && this.currentUser) {
        // Actualizar datos con la información más reciente de Firestore
        this.currentUser = {
          ...this.currentUser,
          ...userData
        };
        console.log('Datos actualizados desde Firestore:', this.currentUser);
      }
    } catch (error) {
      console.warn('No se pudieron recargar los datos desde Firestore:', error);
    }
  }

  // Cerrar sesión
  async onLogout(): Promise<void> {
    try {
      await this.authService.logout();
      this.router.navigate(['/login']);
    } catch (error) {
      console.error('Error al cerrar sesión:', error);
    }
  }

  // Iniciar sesión (redirigir al login)
  onLogin(): void {
    this.router.navigate(['/login']);
  }

  // Obtener opciones por categoría
  getOptionsByCategory(category: string): MenuOption[] {
    return this.menuOptions.filter(option => option.category === category);
  }

  // Navegar a una opción
  navigateToOption(option: MenuOption): void {
    this.router.navigate([option.route]);
  }

  // Obtener todas las categorías únicas
  getCategories(): string[] {
    return [...new Set(this.menuOptions.map(option => option.category))];
  }

  // Obtener nombre de categoría en español
  getCategoryName(category: string): string {
    const categoryNames: {[key: string]: string} = {
      'shopping': 'Compras',
      'profile': 'Mi Perfil',
      'business': 'Mi Negocio'
    };
    return categoryNames[category] || category;
  }

  // Obtener nombre del usuario para mostrar
  getUserDisplayName(): string {
    if (!this.currentUser) return 'Usuario';
    
    // Prioridad: nombre completo > nombre > displayName > email > 'Usuario'
    const nombreCompleto = this.currentUser.nombre && this.currentUser.apellido 
      ? `${this.currentUser.nombre} ${this.currentUser.apellido}`
      : null;
    
    return nombreCompleto ||
           this.currentUser.nombre || 
           this.currentUser.displayName ||
           this.currentUser.email?.split('@')[0] || 
           'Usuario';
  }

  // Obtener rol del usuario con fallback
  getUserRole(): string {
    if (!this.currentUser?.rol) return 'Usuario';
    
    // Mapear roles técnicos a nombres amigables
    const roleNames: {[key: string]: string} = {
      'admin': 'Administrador',
      'vendedor': 'Vendedor',
      'comprador': 'Comprador',
      'usuario': 'Usuario'
    };
    
    return roleNames[this.currentUser.rol] || this.currentUser.rol || 'Usuario';
  }

  // Verificar si el usuario es administrador
  isUserAdmin(): boolean {
    return this.currentUser?.rol === 'admin';
  }

  // Obtener información adicional del usuario
  getUserPhone(): string {
    return this.currentUser?.telefono || 'No disponible';
  }

  getUserAddress(): string {
    return this.currentUser?.direccion || 'No disponible';
  }

  // Verificar si el usuario tiene email verificado
  isEmailVerified(): boolean {
    return this.currentUser?.emailVerified || false;
  }

  // Verificar si el usuario está activo
  isUserActive(): boolean {
    return this.currentUser?.activo !== false; // true por defecto
  }
}
