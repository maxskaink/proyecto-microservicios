import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, OnInit } from '@angular/core';

import { ShoppingCartService } from '../../../service/ShoppinCartService';
import { ItemOrder } from '../../components/item-order/item-order';
import { Order } from '../../../Models/OrderPeticion';
import { Header } from '../../templates/header/header';
import { ArrowLeft } from '../../components/arrow-left/arrow-left';
import { Router } from '@angular/router';
import { map } from 'rxjs';
import { stat } from 'fs';

@Component({
  selector: 'app-history-orders-user',
  imports: [CommonModule, ItemOrder, Header],
  templateUrl: './history-orders-user.html',
  styleUrl: './history-orders-user.css',
})
export class HistoryOrdersUser implements OnInit{ 

  public orders?: Order[] ;
  public oldOrders: Order[] = [];
  public isLoading: boolean = true;
  constructor(
    private shoppingService: ShoppingCartService,
    private cdr: ChangeDetectorRef,
    private router: Router
  ) {

  }
  ngOnInit(): void {
    this.loadUserOrders();
  }

loadUserOrders() {
  this.shoppingService.getUserOrders('').pipe(
    map((orders) => 
      // Si orders es null, lo convertimos a array vacío
      (orders ?? []).map(order => order ?? {}) // cada order null se convierte en objeto vacío
    )
  ).subscribe({
    next: (orders) => {
      this.orders = orders
      console.log("Ordenes: ", this.orders, this.oldOrders)
      this.isLoading = false;
      this.cdr.detectChanges();
    },
    error: (error) => {
      console.error('Error al cargar las órdenes del usuario:', error);
      this.isLoading = false;
    }
  });
}


  goToHome(){
    this.router.navigate(['/home']);
  }
  viewOrderDetails(orderId: string) {
      // Por ahora solo mostramos un alert, después se puede navegar a una página de detalles
      console.log('Ver detalles de la orden:', orderId);
      this.router.navigate(['list-order/view-order', orderId]);
  }

  handlerOrderAction(event: { status: string; id: string }) {
    console.log("estado a enviar desde componete", event.status)
    this.shoppingService.updateStateOrder(event.status, event.id).subscribe({
      next: () => {
        this.loadUserOrders();
      },
      error: (error) => {
        console.error('❌ Error al completar la orden:', error);
      }
    });
  }
}
