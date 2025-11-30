import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, OnChanges, OnInit, Output, SimpleChanges } from '@angular/core';
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
export class FormUser implements  OnInit, OnChanges {
  @Input() userData!: UserResponseBack;
  @Output() submitEvent = new EventEmitter<userPeticion>;
  public form!: FormGroup;
  constructor(
    private fb: FormBuilder,
  ) {
   
  }

 ngOnInit(): void {
    this.initializeForm();
  }

  ngOnChanges(changes: SimpleChanges) {
    if (!this.form) return;
    if (changes['user_data']?.currentValue) {
      this.form.patchValue(this.userData);
    }
  }
  /**
   * Inicializa el formulario con las validaciones necesarias
   */
  private initializeForm(): void {
    this.form = this.fb.group({
      name: [this.userData.name, [Validators.required, Validators.minLength(3), Validators.maxLength(50)]],
      address: [this.userData.profile.address, [Validators.required, Validators.minLength(3), Validators.maxLength(50)]],
      phone: [this.userData.profile.phone, [Validators.required, Validators.maxLength(10), Validators.minLength(10), Validators.pattern('^[0-9]+$')]],
      description: [this.userData.profile.description, [Validators.required, Validators.minLength(10), Validators.maxLength(500)]],
     
    });
  }
  private patchFormWithUser(user: UserResponseBack): void {
  this.form.patchValue({
    name: user.name || '',
    address: user.profile?.address || '',
    phone: user.profile?.phone || '',
    description: user.profile?.description || '',
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
