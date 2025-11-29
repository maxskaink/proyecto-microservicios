import { ChangeDetectorRef, Component, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, NgForm } from '@angular/forms';
import { AuthService } from '../../../service/Authser.vice';
import { Router } from '@angular/router';
import { TenantService } from '../../../service/TenantService';
import { Tenant } from '../../../Models/Tenant';

@Component({
  selector: 'app-login',
  imports: [CommonModule, FormsModule],
  templateUrl: './login.html',
  styleUrl: './login.css',
})
export class Login implements OnInit {
  
  public tenantId: string = '';
  public correo: string = '';
  public password: string = '';
  
  public isLoading: boolean = false;
  public showError: boolean = false;
  public showSuccess: boolean = false;
  public errorType: string = '';
  
  public tenantError: boolean = false;
  public emailError: boolean = false;
  public passwordError: boolean = false;
  
  public tenants = signal<Tenant[]>([]);
  public loadingTenants = signal<boolean>(false);

  constructor(
    private authService: AuthService, 
    private router: Router, 
    private tenantService: TenantService,
    private cdr: ChangeDetectorRef
  ) {}

  ngOnInit(): void {
    this.loadTenants();
  }
  /**
   *  Maneja el proceso de inicio de sesión cuando se envía el formulario.
   * @param loginForm formulario pasado por parametro 
   * @returns promise void 
   */
  async onLogin(loginForm: NgForm): Promise<void> {
    this.hideMessages();
    this.clearFieldErrors();

    if (loginForm.invalid) {
      Object.keys(loginForm.controls).forEach(key => {
        loginForm.controls[key].markAsTouched();
      });
      return;
    }

    this.isLoading = true;

    try {
      const result = await this.authService.login(
        this.correo, 
        this.password, 
        this.tenantId
      );

      this.showSuccess = true;
      
      setTimeout(() => {
        this.router.navigate(['/home']);
      }, 1000);

    } catch (error: any) {
      console.error('Error de login:', error);
      this.handleLoginError(error);
      
    } finally {
      this.isLoading = false;
    }
  }
  /**
   * Carga la lista de tenants desde el servicio TenantService.
   */
  private loadTenants(): void {
    this.loadingTenants.set(true);
    
    this.tenantService.getTenants().subscribe({
      next: (tenants: Tenant[]) => {
        // ✅ Los signals se actualizan automáticamente en el HTML
        this.tenants.set(tenants);
        console.log('Tenants loaded:', tenants);
        this.loadingTenants.set(false);
      },
      error: (error) => {
        console.error('Error loading tenants:', error);
        this.showError = true;
        this.errorType = 'tenant-load';
        this.loadingTenants.set(false);
      }
    });
  }
  /**
   * Maneja los errores de inicio de sesión y actualiza el estado de la interfaz de usuario en consecuencia.
   * @param error objeto de error recibido durante el proceso de inicio de sesión
   */
  private handleLoginError(error: any): void {
    this.showError = true;

    if (error.code === 'auth/invalid-email' || error.message?.includes('correo')) {
      this.errorType = 'email';
      this.emailError = true;
    } else if (error.code === 'auth/wrong-password' || error.message?.includes('contraseña')) {
      this.errorType = 'password';
      this.passwordError = true;
    } else if (error.message?.includes('tenant') || error.message?.includes('tienda')) {
      this.errorType = 'tenant';
      this.tenantError = true;
    } else {
      this.errorType = 'general';
    }

    setTimeout(() => {
      this.hideMessages();
    }, 5000);
  }

  /**
   * Maneja el cambio de selección del tenant.
   */
  public onTenantChange(): void {
    this.tenantError = false;
    console.log('Tenant selected:', this.tenantId);
  }
  /**
   * Maneja el cambio en los campos de entrada del formulario.
   */
  public onInputChange(): void {
    this.clearFieldErrors();
  }
  /**
   * Oculta los mensajes de error y éxito.
   */
  public hideMessages(): void {
    this.showError = false;
    this.showSuccess = false;
    this.errorType = '';
  }
  /**
   * Limpia los errores específicos de los campos del formulario.
   */
  private clearFieldErrors(): void {
    this.tenantError = false;
    this.emailError = false;
    this.passwordError = false;
  }

  /**
   * Inicia sesión con Google
   */
  async loginWithGoogle(): Promise<void> {
    // Verificar que se haya seleccionado un tenant
    if (!this.tenantId) {
      this.tenantError = true;
      this.showError = true;
      this.errorType = 'tenant';
      setTimeout(() => this.hideMessages(), 5000);
      return;
    }

    this.hideMessages();
    this.clearFieldErrors();
    this.isLoading = true;
    this.cdr.detectChanges();
    

    try {
      await this.authService.loginWithGoogle(this.tenantId);
      
      this.showSuccess = true;
      
      setTimeout(() => {
        this.router.navigate(['/home']);
      }, 1000);

    } catch (error: any) {
      console.error('Error de login con Google:', error);
      this.handleLoginError(error);
      
    } finally {
      this.isLoading = false;
      this.cdr.detectChanges();
    }
  }
}