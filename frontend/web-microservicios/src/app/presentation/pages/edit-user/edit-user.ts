import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component } from '@angular/core';
import { Header } from '../../templates/header/header';
import { FormUser } from '../../templates/form-user/form-user';

@Component({
  selector: 'app-edit-user',
  imports: [CommonModule, Header, FormUser],
  templateUrl: './edit-user.html',
  styleUrl: './edit-user.css',
})
export class EditUser {
  stored = localStorage.getItem('user_data');

  constructor(
    private cdr: ChangeDetectorRef
  ){ }

}
