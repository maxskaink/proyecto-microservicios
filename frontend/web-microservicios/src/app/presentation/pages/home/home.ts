import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { Header } from '../../templates/header/header';

@Component({
  selector: 'app-home',
  imports: [CommonModule, Header],
  templateUrl: './home.html',
  styleUrl: './home.css',
})
export class Home {

}
