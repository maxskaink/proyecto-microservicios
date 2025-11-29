import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { User } from '../../pages/user/user';
import { UserResponseBack } from '../../../Models/UserReponseBack';

@Component({
  selector: 'app-user-box',
  imports: [CommonModule],
  standalone : true,
  templateUrl: './user-box.html',
  styleUrl: './user-box.css',
})
export class UserBox {
  @Input() user!: UserResponseBack;
  @Output() action = new EventEmitter<string>();

  clickAction() {
    this.action.emit(this.user.id);
  }

  /**
   * Obtiene la clase CSS para el badge del rol
   */
  getRoleClass(role: string): string {
    switch (role?.toLowerCase()) {
      case 'admin':
        return 'badge-admin';
      case 'cliente':
        return 'badge-client';
      case 'vendedor':
        return 'badge-producer';
      default:
        return 'badge-default';
    }
  }

  /**
   * Obtiene el icono para cada rol
   */
  getRoleIcon(role: string): string {
    switch (role?.toLowerCase()) {
      case 'admin':
        return 'bi bi-shield-check';
      case 'cliente':
        return 'bi bi-person';
      case 'vendedor':
        return 'bi bi-shop';
      default:
        return 'bi bi-person-circle';
    }
  }

  /**
   * Obtiene el label en español para cada rol
   */
  getRoleLabel(role: string): string {
    switch (role?.toLowerCase()) {
      case 'admin':
        return 'Administrador';
      case 'cliente':
        return 'Cliente';
      case 'vendedor':
        return 'Vendedor';
      default:
        return 'Usuario';
    }
  }
}
