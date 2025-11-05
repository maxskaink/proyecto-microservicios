import { Component, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, NgForm } from '@angular/forms';
import { AuthService } from '../../../service /Authser.vice';
import { Router } from '@angular/router';
import { TenantService } from '../../../service /TenantService';
import { Tenant } from '../../../Models/Tenant';
@Component({
  selector: 'app-login',
  imports: [CommonModule, FormsModule],
  templateUrl: './login.html',
  styleUrl: './login.css',
})
export class Login implements OnInit {
  correo: string = '';
  nombre: string = '';
  password: string = '';
  tenantId: string = '';
  tenants: Tenant[] = [];
  
  // Estados para feedback visual
  isLoading: boolean = false;
  showError: boolean = false;
  errorMessage: string = '';
  showSuccess: boolean = false;
  successMessage: string = '';
  errorType: string = '';
  
  // Estados para errores específicos de campos
  emailError: string = '';
  passwordError: string = '';
  tenantError: string = '';
  
  // Estados para modal de credenciales
  showCredentialsModal: boolean = false;
  credentialsMessage: string = '';
  
  // Contador de intentos fallidos
  failedAttempts: number = 0;
  maxAttempts: number = 5;
  isBlocked: boolean = false;
  blockTimeRemaining: number = 0;



  public router = inject(Router);

  constructor(private authService: AuthService, private tenantSerice: TenantService) {}

  ngOnInit(): void {
    this.loadTenants();
  }

  /**
   * Carga la lista de tenants desde el servicio
   */
  private loadTenants(): void {
    this.tenantSerice.getTenants().subscribe({
      next: (tenants) => {
        this.tenants = tenants;
        console.log('Tenants cargados:', tenants);
      },
      error: (error) => {
        console.error('Error al obtener tenants:', error);
        this.showErrorMessage('Error al cargar las tiendas disponibles');
      }
    });
  }

  navigateTo(path: string) {
    this.router.navigate([path]);
  }
  async onLogin(form: NgForm) {
    if (form.invalid) {
      this.showErrorMessage('Por favor, completa todos los campos correctamente');
      return;
    }

    // Validar que se haya seleccionado un tenant
    if (!this.tenantId) {
      this.tenantError = 'Debes seleccionar una tienda';
      this.showErrorMessage('Por favor, selecciona una tienda');
      return;
    }

    if (this.isBlocked) {
      this.showErrorMessage(`Demasiados intentos fallidos. Intenta nuevamente en ${this.blockTimeRemaining} segundos`);
      return;
    }

    this.isLoading = true;
    this.hideMessages();

    try {
      console.log('Starting login process...');
      console.log('Tenant seleccionado:', this.tenantId);
      
      // Pasar el tenantId al método login
      await this.authService.login(this.correo, this.password, this.tenantId);
      this.navigateTo('/home');
    } catch (error: any) {
      this.handleLoginError({
        success: false,
        errorType: 'network_error',
        errorMessage: 'Error de conexión. Verifica tu internet e intenta nuevamente'
      });
      console.error('Error de login:', error);
    } finally {
      this.isLoading = false;
    }
  }
  private getTenantsList() {
    this.tenantSerice.getTenants().subscribe({
      next: (tenants) => {
        this.tenants = tenants;
      },
      error: (error) => {
        console.error('Error al obtener tenants:', error);
      }
    });
  }
  private handleLoginError(result: { success: boolean; errorType?: string; errorMessage?: string; blocked?: boolean; retryAfter?: number }) {
    this.failedAttempts++;
    this.errorType = result.errorType || 'other';
    
    // Limpiar errores previos de campos
    this.clearFieldErrors();
    
    // Si la IP está bloqueada, manejar de forma especial
    if (result.blocked) {
      this.handleIPBlocked(result.retryAfter || 900);
      return;
    }
    
    // Manejar errores específicos de campos
    switch (result.errorType) {
      case 'user_not_found':
        this.emailError = 'El usuario no existe';
        this.showCredentialsModal = true;
        this.credentialsMessage = this.buildCredentialsMessage();
        break;
      case 'wrong_password':
        this.passwordError = 'Contraseña incorrecta';
        this.showCredentialsModal = true;
        this.credentialsMessage = this.buildCredentialsMessage();
        break;
      case 'invalid_email':
        this.emailError = 'El formato del correo electrónico no es válido';
        break;
      case 'invalid_credential':
        this.emailError = 'Correo o contraseña incorrectos';
        this.passwordError = 'Correo o contraseña incorrectos';
        this.showCredentialsModal = true;
        this.credentialsMessage = this.buildCredentialsMessage();
        break;
      case 'account_disabled':
        this.showErrorMessage('Tu cuenta ha sido deshabilitada. Contacta al administrador.');
        break;
      case 'account_locked':
        this.showErrorMessage('Tu cuenta ha sido bloqueada por seguridad. Contacta al administrador.');
        break;
      case 'email_not_verified':
        this.emailError = 'Debes verificar tu correo electrónico';
        break;
      case 'weak_password':
        this.passwordError = 'La contraseña debe tener al menos 6 caracteres';
        break;
      case 'too_many_attempts':
        this.showErrorMessage('Demasiados intentos fallidos. Intenta más tarde.');
        this.blockUser();
        break;
      case 'ip_blocked':
        this.showErrorMessage('Tu IP ha sido bloqueada temporalmente por demasiados intentos fallidos.');
        this.handleIPBlocked(result.retryAfter || 900);
        return;
      case 'network_error':
        this.showErrorMessage('Error de conexión. Verifica tu internet e intenta nuevamente.');
        break;
      case 'server_error':
        this.showErrorMessage('Error del servidor. Intenta nuevamente en unos minutos.');
        break;
      default:
        this.showErrorMessage(result.errorMessage || 'Error de autenticación');
    }
    
    // Bloquear usuario si excede intentos máximos
    if (this.failedAttempts >= this.maxAttempts) {
      this.blockUser();
    }
  }

  private buildCredentialsMessage(): string {
    const remaining = this.maxAttempts - this.failedAttempts;
    let message = `Credenciales incorrectas. Intento ${this.failedAttempts}/${this.maxAttempts}`;
    
    if (remaining > 0) {
      message += ` (${remaining} intentos restantes)`;
    }
    
    message += '\n\nVerifica que:';
    message += '\n• El correo electrónico esté escrito correctamente';
    message += '\n• La contraseña sea la correcta';
    
    return message;
  }

  /**
   * Maneja el cambio de selección del tenant
   */
  onTenantChange(): void {
    this.tenantError = ''; // Limpiar error cuando se selecciona un tenant
    console.log('Tenant seleccionado:', this.tenantId);
  }

  private clearFieldErrors() {
    this.emailError = '';
    this.passwordError = '';
    this.tenantError = '';
  }

  private handleIPBlocked(retryAfter: number) {
    this.isBlocked = true;
    this.blockTimeRemaining = retryAfter;
    
    const errorMsg = `Tu IP ha sido bloqueada temporalmente por demasiados intentos fallidos. Intenta nuevamente en ${this.formatTime(retryAfter)}.`;
    this.showErrorMessage(errorMsg);
    
    const interval = setInterval(() => {
      this.blockTimeRemaining--;
      if (this.blockTimeRemaining <= 0) {
        clearInterval(interval);
        this.isBlocked = false;
        this.failedAttempts = 0;
        this.hideMessages();
      }
    }, 1000);
  }

  private blockUser() {
    this.isBlocked = true;
    this.blockTimeRemaining = 300; // 5 minutos
    
    const interval = setInterval(() => {
      this.blockTimeRemaining--;
      if (this.blockTimeRemaining <= 0) {
        clearInterval(interval);
        this.isBlocked = false;
        this.failedAttempts = 0;
        this.hideMessages();
      }
    }, 1000);
  }

  private showErrorMessage(message: string) {
    this.errorMessage = message;
    this.showError = true;
    this.showSuccess = false;
  }

  private showSuccessMessage(message: string) {
    this.successMessage = message;
    this.showSuccess = true;
    this.showError = false;
  }

  hideMessages() {
    this.showError = false;
    this.showSuccess = false;
  }

  onInputChange() {
    this.hideMessages();
    this.clearFieldErrors();
  }

  closeCredentialsModal() {
    this.showCredentialsModal = false;
    this.credentialsMessage = '';
  }

  formatTime(seconds: number): string {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
  }
}
