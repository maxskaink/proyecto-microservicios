import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ViewUserInfoPage } from './view-user-info-page';

describe('ViewUserInfoPage', () => {
  let component: ViewUserInfoPage;
  let fixture: ComponentFixture<ViewUserInfoPage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ViewUserInfoPage]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ViewUserInfoPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
