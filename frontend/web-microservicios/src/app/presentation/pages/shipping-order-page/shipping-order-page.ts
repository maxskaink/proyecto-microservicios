import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, Input, OnInit } from '@angular/core';
import { ShippingService } from '../../../service/shipping.service';

import { ShoppingCartService } from '../../../service/ShoppinCartService';
import { Order } from '../../../Models/OrderPeticion';
import { shippingItem } from '../../templates/shippingItem/shippingItem';
import { ItemOrder } from '../../components/item-order/item-order';
import { ShippingResponse } from '../../../Models/ShippingResponse';
import { Header } from '../../templates/header/header';
import { ActivatedRoute } from '@angular/router';
import { forkJoin } from 'rxjs';


@Component({
  selector: 'app-shipping-order-page',
  imports: [CommonModule, shippingItem, ItemOrder, Header],
  templateUrl: './shipping-order-page.html',
  styleUrl: './shipping-order-page.css',
})

export class ShippingOrderPage implements OnInit {

  idOrder!: string;
  public shippingInfo!: ShippingResponse;
  public orderInfo!: Order;
  public showActions: boolean = true;
  isLoading = true;

  constructor(
    private route: ActivatedRoute,
    private shpippingOrderService: ShippingService,
    private cdr: ChangeDetectorRef,
    private orderService: ShoppingCartService
  ) {}

  ngOnInit(): void {
    this.idOrder = this.route.snapshot.paramMap.get('id')!;
    this.loadData();
  }
  
  loadData(): void {
    this.isLoading = true;
    forkJoin({
      shipping: this.shpippingOrderService.getShippingByOrderId(this.idOrder),
      order: this.orderService.getOrderById(this.idOrder)
    }).subscribe({
      next: (res) => {
        this.shippingInfo = res.shipping;
        this.orderInfo = res.order;

        this.isLoading = false;
        this.cdr.detectChanges();
      },
      error: (error) => {
        console.error('❌ Error cargando data:', error);
        this.isLoading = false;
      }
    });
  }

  onActionEvent(event: { status: string; id: string }): void {
    console.log('Evento de acción recibido en ShippingOrderPage:', event);

    this.shpippingOrderService.updateStatusShipping(event.id, event.status).subscribe({
      next: (success) => {
        if (success) {
          console.log(`Estado del envío actualizado a "${event.status}" para el ID: ${event.id}`);
        } else {
          console.error(`No se pudo actualizar el estado del envío para el ID: ${event.id}`);
        }
      },
      error: (error) => {
        console.error(`Error al actualizar el estado del envío para el ID: ${event.id}`, error);
      }
    });
  }

}
