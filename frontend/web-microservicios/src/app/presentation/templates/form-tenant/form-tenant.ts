import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Output } from '@angular/core';
import { TenantPeticion } from '../../../Models/TenantPeticion';
import { FormBuilder, FormGroup, NgForm, ReactiveFormsModule, Validators } from '@angular/forms';

@Component({
  selector: 'app-form-tenant',
  imports: [CommonModule, ReactiveFormsModule],
  standalone: true,
  templateUrl: './form-tenant.html',
  styleUrl: './form-tenant.css',
})
export class FormTenant {
  @Output() tenant = new EventEmitter<TenantPeticion>();;
  public tenantForm!: FormGroup;

  constructor(
    private fb: FormBuilder,
  ) {
    this.initializeForm();
  }
  /**
   * emite el tenatn creado
   */
  onSubmit() {
    const tenantPeticion = this.tenantForm.value as TenantPeticion;
    this.tenant.emit(tenantPeticion);
  }
  private initializeForm(): void {
    this.tenantForm = this.fb.group({
      tenant_id: ['', [Validators.required, Validators.pattern(/^[a-zA-Z0-9-]+$/), Validators.minLength(3), Validators.maxLength(20)]],
      tenant_name: ['', [Validators.required, Validators.minLength(3), Validators.maxLength(50)]],
      location: ['', [Validators.required, Validators.minLength(3), Validators.maxLength(50)]],
      description: ['', [Validators.required, Validators.minLength(10) , Validators.maxLength(200)]],
    });
  }
  
}
