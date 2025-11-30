import { ChangeDetectorRef, Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { CommonModule } from '@angular/common';
import { UserResponseBack } from '../../../Models/UserReponseBack';
import { AuthService } from '../../../service/Authser.vice';
import { Observable } from 'rxjs';
import { BlobOptions } from 'buffer';

@Component({
  selector: 'app-view-info-producer',
  imports: [CommonModule],
  templateUrl: './view-info-producer.html',
  styleUrl: './view-info-producer.css',
})
export class ViewInfoProducer implements OnInit {
  
  @Input() producer?: UserResponseBack;
  @Input() showActions: boolean = false;
  @Input() isEditable: boolean = false;
  @Output() action = new EventEmitter<string>();
  isAdmin: boolean = false;

  clickAction() {
    this.action.emit('some-action');
  }
  constructor(
    private cdr: ChangeDetectorRef,
  ) {}
  ngOnInit(): void {
    this.consultRol();
  }

  consultRol() {
     const user = JSON.parse(localStorage.getItem('user_data') || 'null');
     if (user && user.rol === 'admin') {
       this.isAdmin = true;
       this.cdr.detectChanges();
     }
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
