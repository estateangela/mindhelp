import 'package:json_annotation/json_annotation.dart';

part 'counselor.g.dart';

@JsonSerializable()
class Counselor {
  final String id;
  final String name;
  final String licenseNumber;
  final String gender;
  final String specialties;
  final String workLocation;
  final String workUnit;

  Counselor({
    required this.id,
    required this.name,
    required this.licenseNumber,
    required this.gender,
    required this.specialties,
    required this.workLocation,
    required this.workUnit,
  });

  factory Counselor.fromJson(Map<String, dynamic> json) => _$CounselorFromJson(json);
  Map<String, dynamic> toJson() => _$CounselorToJson(this);
}

@JsonSerializable()
class CounselorListResponse {
  final List<Counselor> counselors;
  final int total;
  final int page;
  final int pageSize;

  CounselorListResponse({
    required this.counselors,
    required this.total,
    required this.page,
    required this.pageSize,
  });

  factory CounselorListResponse.fromJson(Map<String, dynamic> json) => 
      _$CounselorListResponseFromJson(json);
  Map<String, dynamic> toJson() => _$CounselorListResponseToJson(this);
}

@JsonSerializable()
class CounselingCenter {
  final String id;
  final String name;
  final String address;
  final String phone;
  final String? website;
  final String? description;
  final bool hasOnlineService;

  CounselingCenter({
    required this.id,
    required this.name,
    required this.address,
    required this.phone,
    this.website,
    this.description,
    required this.hasOnlineService,
  });

  factory CounselingCenter.fromJson(Map<String, dynamic> json) => 
      _$CounselingCenterFromJson(json);
  Map<String, dynamic> toJson() => _$CounselingCenterToJson(this);
}

@JsonSerializable()
class CounselingCenterListResponse {
  final List<CounselingCenter> counselingCenters;
  final int total;
  final int page;
  final int pageSize;

  CounselingCenterListResponse({
    required this.counselingCenters,
    required this.total,
    required this.page,
    required this.pageSize,
  });

  factory CounselingCenterListResponse.fromJson(Map<String, dynamic> json) => 
      _$CounselingCenterListResponseFromJson(json);
  Map<String, dynamic> toJson() => _$CounselingCenterListResponseToJson(this);
}

@JsonSerializable()
class RecommendedDoctor {
  final String id;
  final String name;
  final String specialty;
  final String location;
  final String phone;
  final String? description;
  final int experienceCount;

  RecommendedDoctor({
    required this.id,
    required this.name,
    required this.specialty,
    required this.location,
    required this.phone,
    this.description,
    required this.experienceCount,
  });

  factory RecommendedDoctor.fromJson(Map<String, dynamic> json) => 
      _$RecommendedDoctorFromJson(json);
  Map<String, dynamic> toJson() => _$RecommendedDoctorToJson(this);
}

@JsonSerializable()
class RecommendedDoctorListResponse {
  final List<RecommendedDoctor> recommendedDoctors;
  final int total;
  final int page;
  final int pageSize;

  RecommendedDoctorListResponse({
    required this.recommendedDoctors,
    required this.total,
    required this.page,
    required this.pageSize,
  });

  factory RecommendedDoctorListResponse.fromJson(Map<String, dynamic> json) => 
      _$RecommendedDoctorListResponseFromJson(json);
  Map<String, dynamic> toJson() => _$RecommendedDoctorListResponseToJson(this);
}
