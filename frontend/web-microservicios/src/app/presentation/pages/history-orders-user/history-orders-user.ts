import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, OnInit } from '@angular/core';

import { ShoppingCartService } from '../../../service/ShoppinCartService';
import { ItemOrder } from '../../components/item-order/item-order';
import { Order } from '../../../Models/OrderPeticion';
import { Header } from '../../templates/header/header';
import { ArrowLeft } from '../../components/arrow-left/arrow-left';
import { Router } from '@angular/router';

@Component({
  selector: 'app-history-orders-user',
  imports: [CommonModule, ItemOrder, Header, ArrowLeft],
  templateUrl: './history-orders-user.html',
  styleUrl: './history-orders-user.css',
})
export class HistoryOrdersUser implements OnInit{ 

  public orders: Order[] = [];
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
    this.shoppingService.getUserOrders('').subscribe({
      next: (orders) => {
        console.log('Órdenes del usuario:', orders);
        this.orders = orders;
        this.isLoading = false;
        this.cdr.detectChanges();
        
      },
      error: (error) => {
        console.error('Error al cargar las órdenes del usuario:', error);
      }
    });
    this.shoppingService.getUserOrders('paid').subscribe({
      next: (ordersOld) => {
        console.log('Órdenes antiguas del usuario:', ordersOld);
        this.oldOrders = ordersOld;
        this.isLoading = false;
        this.cdr.detectChanges();
        
      },
      error: (error) => {
        console.error('Error al cargar las órdenes antiguas del usuario:', error);
      }
    });
  }

  goToHome(){
    this.router.navigate(['/home']);
  }
}
