import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { TenantService } from '../../../service/TenantService';
import { TenantPeticion } from '../../../Models/TenantPeticion';
import { Header } from '../../templates/header/header';
import { FormTenant } from '../../templates/form-tenant/form-tenant';
import Swal from 'sweetalert2';

@Component({
  selector: 'app-register-tenant',
  imports: [CommonModule, ReactiveFormsModule, Header, FormTenant],
  templateUrl: './register-tenant.html',
  styleUrl: './register-tenant.css',
})
export class RegisterTenant implements OnInit {
  
  tenantForm!: FormGroup;
  isSubmitting = false;
  errorMessage = '';
  successMessage = '';

  constructor(
    private fb: FormBuilder,
    private serviceTenant: TenantService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.initializeForm();
  }

  private initializeForm(): void {
    this.tenantForm = this.fb.group({
      tenant_id: ['', [Validators.required, Validators.pattern(/^[a-zA-Z0-9-]+$/)]],
      tenant_name: ['', [Validators.required, Validators.minLength(3)]]
    });
  }

  onRegisterTenant(tenatn: TenantPeticion): void {
   this.isSubmitting = true;
    this.errorMessage = '';
    this.successMessage = '';
    this.serviceTenant.PostTenant(tenatn).subscribe({
      next: (response) => {
        this.isSubmitting = false;
        this.successMessage = 'Tenant registrado exitosamente.';
        this.tenantForm.reset();
        setTimeout(() => {
          this.router.navigate(['/home']);
        }, 2000);
      },
      error: (error) => {
        this.isSubmitting = false;
        this.errorMessage = 'Error al registrar el tenant. Inténtalo de nuevo.';
        console.error('Error al registrar el tenant:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Error al registrar el tenant. Inténtalo de nuevo.',
          buttonsStyling: false,
          customClass: {
            confirmButton: 'btn btn-danger'
          }
        });   
      }
    });
  }
}