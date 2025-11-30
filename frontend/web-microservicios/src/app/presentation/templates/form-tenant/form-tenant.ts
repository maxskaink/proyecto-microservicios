import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Output } from '@angular/core';
import { TenantPeticion } from '../../../Models/TenantPeticion';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators, AbstractControl } from '@angular/forms';

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
      tenant_id: ['', [
        Validators.required, 
        Validators.pattern(/^[a-z0-9-]+$/), // Solo minúsculas, números y guiones
        Validators.minLength(3), 
        Validators.maxLength(20),
        this.noSpacesValidator,
        this.lowercaseValidator
      ]],
      tenant_name: ['', [Validators.required, Validators.minLength(3), Validators.maxLength(50)]],
      location: ['', [Validators.required, Validators.minLength(3), Validators.maxLength(50)]],
      description: ['', [Validators.required, Validators.minLength(10) , Validators.maxLength(100)]],
    });
  }

  /**
   * Validador personalizado para asegurar que no hay espacios
   */
  private noSpacesValidator(control: AbstractControl) {
    if (!control.value) return null;
    if (control.value.includes(' ')) {
      return { hasSpaces: true };
    }
    return null;
  }

  /**
   * Validador personalizado para asegurar que esté en minúsculas
   */
  private lowercaseValidator(control: AbstractControl) {
    if (!control.value) return null;
    if (control.value !== control.value.toLowerCase()) {
      return { notLowercase: true };
    }
    return null;
  }

  /**
   * Convierte automáticamente a minúsculas y remueve espacios
   */
  onTenantIdChange(event: any) {
    let value = event.target.value;
    // Convertir a minúsculas y reemplazar espacios con guiones
    value = value.toLowerCase().replace(/\s+/g, '-');
    // Remover caracteres no permitidos excepto letras, números y guiones
    value = value.replace(/[^a-z0-9-]/g, '');
    
    // Actualizar el valor del control
    this.tenantForm.get('tenant_id')?.setValue(value);
  }
  
}
