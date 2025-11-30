import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, OnChanges, OnInit, SimpleChanges } from '@angular/core';
import { UserBox } from '../../components/user-box/user-box';
import { UserResponseBack } from '../../../Models/UserReponseBack';
import { Header } from '../../templates/header/header';
import { UsersService } from '../../../service/users.service';
import Swal from 'sweetalert2';
import { IsLoading } from '../../components/is-loading/is-loading';
import { LoadingService } from '../../../service/loading-service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-view-panel-users',
  imports: [CommonModule, UserBox, Header, IsLoading],
  templateUrl: './view-panel-users.html',
  styleUrl: './view-panel-users.css',
})
export class ViewPanelUsers implements OnInit{
  public users!: UserResponseBack[];
  public clients!: UserResponseBack[];
  public producers!: UserResponseBack[];
  public admins!: UserResponseBack[];
  constructor(
    private userService: UsersService,
    private cdr:  ChangeDetectorRef,
    private router: Router,
    private isLoading: LoadingService
  ) {}
  ngOnInit(): void {
    this.loadUsers();
  }

  /**
   * cargar usuarios desde el servicio, todos los usuarios existentes del tenant
   */
  loadUsers() {
    this.isLoading.show('Cargando usuarios...');
    this.userService.getAllUsers().subscribe({
    next: (data: UserResponseBack[]) => {
      this.users = data;
      this.clients = this.users.filter(user => user.rol === 'client');
      this.producers = this.users.filter(user => user.rol === 'producer');
      this.admins = this.users.filter(user => user.rol === 'admin');
      this.cdr.detectChanges();
      this.isLoading.hide();
    },
    error: (error) => {
      this.cdr.detectChanges();
      this.isLoading.hide();
      console.error('Error fetching users:', error);
      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: 'There was an error fetching the users.',
        buttonsStyling: false,
        customClass: {
          confirmButton: 'btn btn-primary',
        },
      });
    }
  });
}
    onAction(id: string){
      console.log('Acción recibida para el usuario con ID:', id);
      this.router.navigate(['/user/view-user-info', id]); 
    }


} 
