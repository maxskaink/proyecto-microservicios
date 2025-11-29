import { Component, EventEmitter, Input, Output } from '@angular/core';
import { CommonModule } from '@angular/common';
import { UserResponseBack } from '../../../Models/UserReponseBack';

@Component({
  selector: 'app-view-info-producer',
  imports: [CommonModule],
  templateUrl: './view-info-producer.html',
  styleUrl: './view-info-producer.css',
})
export class ViewInfoProducer {
  @Input() producer?: UserResponseBack;
  @Input() showActions: boolean = false;
  @Output() action = new EventEmitter<string>();

  clickAction() {
    // Implementar según necesidad
  }

  /**
   * Formatea una fecha para mostrar de manera amigable
   */
  formatDate(dateString: string): string {
    if (!dateString) return 'No disponible';
    
    try {
      const date = new Date(dateString);
      const now = new Date();
      const diffTime = Math.abs(now.getTime() - date.getTime());
      const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
      
      if (diffDays < 30) {
        return `Hace ${diffDays} día${diffDays !== 1 ? 's' : ''}`;
      } else if (diffDays < 365) {
        const months = Math.floor(diffDays / 30);
        return `Hace ${months} mes${months !== 1 ? 'es' : ''}`;
      } else {
        const years = Math.floor(diffDays / 365);
        return `Hace ${years} año${years !== 1 ? 's' : ''}`;
      }
    } catch (error) {
      return 'Fecha no válida';
    }
  }
}
