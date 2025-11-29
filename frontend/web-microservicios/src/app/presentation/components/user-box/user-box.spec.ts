import { ComponentFixture, TestBed } from '@angular/core/testing';

import { UserBox } from './user-box';

describe('UserBox', () => {
  let component: UserBox;
  let fixture: ComponentFixture<UserBox>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [UserBox]
    })
    .compileComponents();

    fixture = TestBed.createComponent(UserBox);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
