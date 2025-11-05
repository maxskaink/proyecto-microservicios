import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { TenantService } from '../../../service /TenantService';
import { TenantPeticion } from '../../../Models/TenantPeticion';
import { Header } from '../../templates/header/header';

@Component({
  selector: 'app-register-tenant',
  imports: [CommonModule, ReactiveFormsModule, Header],
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

  onSubmit(): void {
    if (this.tenantForm.invalid) {
      this.tenantForm.markAllAsTouched();
      return;
    }

    this.isSubmitting = true;
    this.errorMessage = '';
    this.successMessage = '';

    const tenantData: TenantPeticion = this.tenantForm.value;

    this.serviceTenant.PostTenant(tenantData).subscribe({
      next: (response) => {
        console.log('✅ Tenant registrado exitosamente:', response);
        this.successMessage = 'Tenant registrado exitosamente';
        this.isSubmitting = false;
        
        setTimeout(() => {
          this.router.navigate(['/user']);
        }, 2000);
      },
      error: (error) => {
        console.error('❌ Error al registrar tenant:', error);
        this.errorMessage = 'Error al registrar el tenant. Intenta de nuevo.';
        this.isSubmitting = false;
      }
    });
  }
}