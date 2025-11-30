import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { User } from '../../pages/user/user';
import { UserResponseBack } from '../../../Models/UserReponseBack';
import { userPeticion } from '../../../Models/UserPeticion';

@Component({
  selector: 'app-form-user',
  imports: [CommonModule, ReactiveFormsModule],
  standalone  : true,
  templateUrl: './form-user.html',
  styleUrl: './form-user.css',
})
export class FormUser {
  @Input() userData!: UserResponseBack;
  @Output() submitEvent = new EventEmitter<userPeticion>;
  public form!: FormGroup;
  constructor(
    private fb: FormBuilder,
  ) {
    this.initializeForm();
  }
  /**
   * Inicializa el formulario con las validaciones necesarias
   */
  private initializeForm(): void {
    this.form = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(3), Validators.maxLength(50)]],
      address: ['', [Validators.required, Validators.minLength(3), Validators.maxLength(50)]],
      phone: ['', [Validators.required, Validators.minLength(10) , Validators.maxLength(200)]],
      description: ['', [Validators.required, Validators.minLength(10), Validators.maxLength(500)]],
    });
  }

  /**
   * Maneja el evento de envío del formulario
   */
  onsubmit(): void {
    if (this.form.valid) {
      const formValue = this.form.value;
      const userPeticionData: userPeticion = {
        email: this.userData.email,
        name: formValue.name,
        profile: {
          address: formValue.address,
          phone: formValue.phone,
          avatar_url: this.userData.profile.avatar_url,
          description: formValue.description,
        },
      };
      this.submitEvent.emit(userPeticionData);
    } else {
      this.form.markAllAsTouched();
    }
  }
}
