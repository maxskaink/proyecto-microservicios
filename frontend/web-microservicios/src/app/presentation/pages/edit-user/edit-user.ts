import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { Header } from '../../templates/header/header';
import { FormUser } from '../../templates/form-user/form-user';
import { IsLoading } from '../../components/is-loading/is-loading';
import { LoadingService } from '../../../service/loading-service';
import { UsersService } from '../../../service/users.service';
import { UserResponseBack } from '../../../Models/UserReponseBack';
import { userPeticion } from '../../../Models/UserPeticion';
import Swal from 'sweetalert2';

@Component({
  selector: 'app-edit-user',
  imports: [CommonModule, Header, FormUser, IsLoading],
  templateUrl: './edit-user.html',
  styleUrl: './edit-user.css',
})
export class EditUser {
  user: UserResponseBack | null = null;

  constructor(
    private cdr: ChangeDetectorRef,
    private loadingService: LoadingService,
    private userService: UsersService
  ) {}

  ngOnInit(): void {
    this.laodUser();
  }
  laodUser(){
    this.loadingService.show("Cargando usuario...");
    const stored = localStorage.getItem('user_data');

    if (stored) {
      try {
        this.user = JSON.parse(stored) as UserResponseBack;
        this.cdr.detectChanges();
        this.loadingService.hide();
      } catch (error) {
        console.error('Error parsing user_data:', error);
        this.user = null;
        this.cdr.detectChanges();
        this.loadingService.hide();
      }
    }
  }
  onSubmit(userData: userPeticion): void {
    this.userService.updateUser(this.user!.id, userData).subscribe({
      next: (updatedUser: UserResponseBack) => {
        console.log('Usuario actualizado:', updatedUser);
        localStorage.setItem('user_data', JSON.stringify(updatedUser));
        this.user = updatedUser;
        this.cdr.detectChanges();
        Swal.fire({
          icon: 'success',
          title: 'Perfil actualizado',
          text: 'Tu información ha sido actualizada correctamente.',
          buttonsStyling: false,
          customClass: {
            confirmButton: 'btn btn-primary',
          },  
        });
      },
      error: (error) => {
        console.error('Error al actualizar el usuario:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Hubo un problema al actualizar tu información. Por favor, intenta de nuevo más tarde.',
          buttonsStyling: false,
          customClass: {
            confirmButton: 'btn btn-primary',
          },
        })
      }
    });   
  }

}
