import { CommonModule } from '@angular/common';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { Observable } from 'rxjs';
import { LoadingService } from '../../../service/loading-service';

@Component({
  selector: 'app-is-loading',
  imports: [CommonModule],
  templateUrl: './is-loading.html',
  styleUrl: './is-loading.css',
})
export class IsLoading implements OnInit, OnDestroy {
  isLoading$: Observable<boolean>;
  loadingMessage$: Observable<string>;

  constructor(private loadingService: LoadingService) {
    this.isLoading$ = this.loadingService.isLoading$;
    this.loadingMessage$ = this.loadingService.loadingMessage$;
  }

  ngOnInit(): void {
    // Componente inicializado
  }

  ngOnDestroy(): void {
    // Cleanup si es necesario
  }
}
