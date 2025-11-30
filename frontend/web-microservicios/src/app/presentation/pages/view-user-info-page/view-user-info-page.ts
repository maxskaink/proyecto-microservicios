import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component } from '@angular/core';
import { ViewInfoProducer } from '../../templates/view-info-producer/view-info-producer';
import { UserResponseBack } from '../../../Models/UserReponseBack';
import { UsersService } from '../../../service/users.service';
import { ActivatedRoute } from '@angular/router';
import { Header } from '../../templates/header/header';
import { LoadingService } from '../../../service/loading-service';
import { IsLoading } from '../../components/is-loading/is-loading';

@Component({
  selector: 'app-view-user-info-page',
  imports: [CommonModule, ViewInfoProducer, Header, IsLoading],
  templateUrl: './view-user-info-page.html',
  styleUrl: './view-user-info-page.css',
})
export class ViewUserInfoPage {
  public user!: UserResponseBack;
  constructor(
    private userService: UsersService,
    private route: ActivatedRoute,
    private cdr:  ChangeDetectorRef,
    private isLoading: LoadingService
  ) {}
  public ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      const userId = params.get('id');
      if (userId) {
        this.loadUser(userId);
      }
    });
  }
  public loadUser(userId: string): void {
    this.isLoading.show('Cargando información del usuario...');
    this.userService.getUserByID(userId).subscribe({
      next: (userData: UserResponseBack) => {
        this.user = userData;
        this.cdr.detectChanges(); 
        this.isLoading.hide();  
      },
      error: (error) => {
        console.error('Error fetching user data:', error);
        this.isLoading.hide();
      },
    });
  }  
  
  handleAction(state:string){
    this.userService.updateRoleUser(this.user.id, state).subscribe({
      next: (response) => {
        console.log('User role updated successfully:', response);
        this.loadUser(this.user.id);
        this.cdr.detectChanges(); 
      },
      error: (error) => {
        console.error('Error updating user role:', error);
      },
    });
  }
}
