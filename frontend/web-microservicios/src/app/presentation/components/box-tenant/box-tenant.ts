import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output, OnInit } from '@angular/core';
import { IsLoading } from '../is-loading/is-loading';
import { Tenant } from '../../../Models/Tenant';


@Component({
  selector: 'app-box-tenant',
  imports: [CommonModule],
  standalone: true,
  templateUrl: './box-tenant.html',
  styleUrl: './box-tenant.css',
})
export class BoxTenant implements OnInit {
  @Input() tenat!: Tenant;
  @Output() action = new EventEmitter<string>();
  
  public images: string[] = [
    'https://chproyectosinmobiliarios.org/wp-content/uploads/2021/07/DSC_2917-1024x678.jpg',
    'https://www.minutaagropecuaria.com/wp-content/uploads/2020/09/D%C3%ADa-mundial-de-la-Agricultura-1.jpg',
    'https://tuvivienda.co/wp-content/uploads/2024/02/1-57.jpg',
  ];
  
  public randomImage: string = '';

  constructor() { }

  ngOnInit(): void {
    this.randomImage = this.getRandomImage();
  }

  /**
   * Obtiene una imagen aleatoria del array de imágenes
   */
  getRandomImage(): string {
    const randomIndex = Math.floor(Math.random() * this.images.length);
    return this.images[randomIndex];
  }

  clickAction() {
    this.action.emit(this.tenat.tenant_id);
  }
}
