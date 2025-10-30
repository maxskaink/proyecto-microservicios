import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule, NgForm } from '@angular/forms';
import { Router } from '@angular/router';
import { SimpleAuthService } from '@org/shared';

@Component({
  selector: 'app-login',
  imports: [CommonModule, FormsModule],
  templateUrl: './login.html',
  styleUrl: './login.css',
})
export class Login {
  correo: string = '';
  nombre: string = '';
  password: string = '';
  
  emailError: boolean = false;
  passwordError: boolean = false;
  isLoading: boolean = false;
  errorMessage: string = '';

  constructor(
    private router: Router,
    private authService: SimpleAuthService
  ) {
    console.log('LoginPage - Constructor iniciado');
  }

  onInputChange() {
    this.emailError = false;
    this.passwordError = false;
    this.errorMessage = '';
  }

  async onLogin(form: NgForm) {
    if (form.invalid) {
      return;
    }
    
    try {
      this.isLoading = true;
      this.emailError = false;
      this.passwordError = false;
      this.errorMessage = '';
  
      console.log('Iniciando login con Firebase:', this.correo);
      const user = await this.authService.login(this.correo, this.password);
      
      console.log('Login exitoso:', user);
      alert(`¡Bienvenido ${user.email}!`);
      this.navigateTo('/home');
      
    } catch (error: any) {
      console.error('Error de login:', error);
      this.errorMessage = error.message;
      if (error.message.includes('email') || error.message.includes('Usuario no encontrado')) {
        this.emailError = true;
      } else if (error.message.includes('contraseña') || error.message.includes('Credenciales')) {
        this.passwordError = true;
      } else {
        this.emailError = true;
        this.passwordError = true;
      }
      
    } finally {
      this.isLoading = false;
    }
  }

  navigateTo(path: string) {
    this.router.navigate([path]);
  }
}