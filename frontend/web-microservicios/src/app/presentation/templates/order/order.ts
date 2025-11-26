import { Component, EventEmitter, Input, OnChanges, Output, SimpleChanges } from '@angular/core';
import { OrderResponse } from '../../../Models/OrderPeticion';
import { BehaviorSubject } from 'rxjs';
import { CommonModule } from '@angular/common';
import { ProductViewBuy } from '../product-view-buy/product-view-buy';
import { ProductBox } from '../../components/product-box/product-box';
import { ProductService } from '../../../service/ProductService';

@Component({
  selector: 'app-order',
  imports: [CommonModule, ProductBox],
  templateUrl: './order.html',
  standalone: true,
  styleUrl: './order.css',
})
export class Order implements OnChanges{
  
  @Input() createdOrderResponse: OrderResponse | null = null;
  @Output() close = new EventEmitter<boolean>();

  constructor(private Service: ProductService) {}

  ngOnChanges(changes: SimpleChanges): void {
    console.log("cambios en order", changes);
  }
  
  closeOrder(){
    this.close.next(true);
  }
  formatPrice(price: number): string {
    return price.toLocaleString('es-CO', {
      style: 'currency',
      currency: 'COP'
    });
  }

}
