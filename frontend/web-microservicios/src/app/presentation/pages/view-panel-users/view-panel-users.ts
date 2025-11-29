import { CommonModule } from '@angular/common';
import { Component, OnChanges, SimpleChanges } from '@angular/core';
import { UserBox } from '../../components/user-box/user-box';
import { UserResponseBack } from '../../../Models/UserReponseBack';
import { Header } from '../../templates/header/header';
import { UsersService } from '../../../service/users.service';
import Swal from 'sweetalert2';

@Component({
  selector: 'app-view-panel-users',
  imports: [CommonModule, UserBox, Header],
  templateUrl: './view-panel-users.html',
  styleUrl: './view-panel-users.css',
})
export class ViewPanelUsers implements OnChanges{
  public users!: UserResponseBack[];
  public producers!: UserResponseBack[];
  public admins!: UserResponseBack[];
  constructor(
    private userService: UsersService,
  ) {}
  ngOnChanges(changes: SimpleChanges): void {
    this.loadUsers();
  }
  /**
   * cargar usuarios desde el servicio, todos los usuarios existentes del tenant
   */
  loadUsers() {
  this.userService.getAllUsers().subscribe({
    next: (data: UserResponseBack[]) => {
      this.users = data.filter(user => user.rol !== 'client');
      this.producers = this.users.filter(user => user.rol === 'producer');
      this.admins = this.users.filter(user => user.rol === 'admin');
    },
    error: (error) => {
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


} 
