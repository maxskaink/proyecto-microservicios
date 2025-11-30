import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { TenantService } from '../../../service/TenantService';
import { TenantPeticion } from '../../../Models/TenantPeticion';
import { Header } from '../../templates/header/header';
import { FormTenant } from '../../templates/form-tenant/form-tenant';
import Swal from 'sweetalert2';
import { LoadingService } from '../../../service/loading-service';
import { AuthService } from '../../../service/Authser.vice';
import { firstValueFrom } from 'rxjs';

@Component({
  selector: 'app-register-tenant',
  imports: [CommonModule, ReactiveFormsModule, Header, FormTenant],
  templateUrl: './register-tenant.html',
  styleUrl: './register-tenant.css',
})
export class RegisterTenant  {
  
  
  isSubmitting = false;
  errorMessage = '';
  successMessage = '';

  constructor(
    private fb: FormBuilder,
    private serviceTenant: TenantService,
    private router: Router,
    private loadingService: LoadingService,
    private authService: AuthService,
  ) {}

 

  async onRegisterTenant(tenant: TenantPeticion): Promise<void> {
    this.isSubmitting = true;
    this.loadingService.show('Registrando tenant...');
    this.errorMessage = '';
    this.successMessage = '';

    try {
      // 1️⃣ Crear tenant
      const response = await firstValueFrom(this.serviceTenant.PostTenant(tenant));

      this.successMessage = 'Tenant registrado exitosamente.';
      Swal.fire({
        icon: 'success',
        title: 'Éxito',
        text: 'Tenant registrado exitosamente.',
        buttonsStyling: false,
        customClass: { confirmButton: 'btn btn-success' }
      });

      // 2️⃣ Ahora SÍ logear usando el tenant creado
      console.log('Iniciando sesión con el tenant creado:', tenant.tenant_id);
      await this.authService.fetchUserFromBackend(tenant.tenant_id);

      // 3️⃣ Redirigir al home solo cuando termine
      this.router.navigate(['/home']);

    } catch (error) {
      console.error('Error al registrar tenant o iniciar sesión:', error);
      this.errorMessage = 'Error al registrar el tenant. Inténtalo de nuevo.';

      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: 'Error al registrar el tenant. Inténtalo de nuevo.',
        buttonsStyling: false,
        customClass: { confirmButton: 'btn btn-danger' }
      });
    } finally {
      this.isSubmitting = false;
      this.loadingService.hide();
    }
  }

}