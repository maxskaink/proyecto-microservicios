import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Firebase } from './firebase';

describe('Firebase', () => {
  let component: Firebase;
  let fixture: ComponentFixture<Firebase>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Firebase],
    }).compileComponents();

    fixture = TestBed.createComponent(Firebase);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
