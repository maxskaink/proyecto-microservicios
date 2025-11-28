import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class LoadingService {
  private _isLoading$ = new BehaviorSubject<boolean>(false);
  private _loadingMessage$ = new BehaviorSubject<string>('Cargando...');
  
  // Contador para manejar múltiples operaciones simultáneas
  private loadingCounter = 0;

  constructor() {}

  /**
   * Observable para saber si está cargando
   */
  get isLoading$(): Observable<boolean> {
    return this._isLoading$.asObservable();
  }

  /**
   * Observable para el mensaje de carga
   */
  get loadingMessage$(): Observable<string> {
    return this._loadingMessage$.asObservable();
  }

  /**
   * Obtiene el estado actual de loading
   */
  get isLoading(): boolean {
    return this._isLoading$.value;
  }

  /**
   * Muestra la pantalla de carga con mensaje opcional
   * @param message - Mensaje personalizado para mostrar
   */
  show(message: string = 'Cargando...'): void {
    this.loadingCounter++;
    this._loadingMessage$.next(message);
    this._isLoading$.next(true);
  }

  /**
   * Oculta la pantalla de carga
   */
  hide(): void {
    this.loadingCounter--;
    
    // Solo ocultar cuando no hay más operaciones pendientes
    if (this.loadingCounter <= 0) {
      this.loadingCounter = 0;
      this._isLoading$.next(false);
    }
  }

  /**
   * Fuerza el ocultado de la pantalla de carga
   * Útil para casos de error donde se necesita limpiar el estado
   */
  forceHide(): void {
    this.loadingCounter = 0;
    this._isLoading$.next(false);
  }

  /**
   * Actualiza solo el mensaje sin cambiar el estado de loading
   * @param message - Nuevo mensaje a mostrar
   */
  updateMessage(message: string): void {
    this._loadingMessage$.next(message);
  }

  /**
   * Ejecuta una operación asíncrona mostrando loading automáticamente
   * @param operation - Función que retorna una Promise
   * @param message - Mensaje a mostrar durante la carga
   */
  async withLoading<T>(
    operation: () => Promise<T>,
    message: string = 'Procesando...'
  ): Promise<T> {
    this.show(message);
    try {
      const result = await operation();
      this.hide();
      return result;
    } catch (error) {
      this.hide();
      throw error;
    }
  }

  /**
   * Versión para Observables
   * @param operation - Función que retorna un Observable
   * @param message - Mensaje a mostrar durante la carga
   */
  withLoadingObservable<T>(
    operation: () => Observable<T>,
    message: string = 'Procesando...'
  ): Observable<T> {
    this.show(message);
    return new Observable<T>(subscriber => {
      const subscription = operation().subscribe({
        next: (value) => subscriber.next(value),
        error: (error) => {
          this.hide();
          subscriber.error(error);
        },
        complete: () => {
          this.hide();
          subscriber.complete();
        }
      });

      // Cleanup cuando se cancela la suscripción
      return () => {
        this.hide();
        subscription.unsubscribe();
      };
    });
  }
}
